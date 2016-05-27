package commands

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/app-restarter/api"
	"github.com/cloudfoundry-incubator/app-restarter/commands/displayhelpers"
	"github.com/cloudfoundry-incubator/app-restarter/models"
	"github.com/cloudfoundry-incubator/app-restarter/resource_mapper"
	"github.com/cloudfoundry-incubator/app-restarter/ui"
	"sync"
)

const (
	Success = iota
	Stopped
	Warning
	Err
)

type RestartAppsExecutor struct {
	AppsGetterFunc resource_mapper.AppsGetterFunc
	RestartAppsUI  *ui.RestartApps
}

func (exe *RestartAppsExecutor) Execute(cliConnection api.Connection) error {
	exe.RestartAppsUI.BeforeAll() //move me to the command

	apiClient, err := api.NewClient(cliConnection)
	if err != nil {
		return err
	}

	appRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetAppsRequest),
	)

	appPaginatedRequester, err := api.NewPaginatedRequester(cliConnection, appRequestFactory)
	if err != nil {
		return err
	}

	apps, err := exe.AppsGetterFunc(
		models.ApplicationsParser{},
		appPaginatedRequester,
	)
	if err != nil {
		return err
	}

	spaceRequestFactory := apiClient.HandleFiltersAndParameters(
		apiClient.Authorize(apiClient.NewGetSpacesRequest),
	)

	spacePaginatedRequester, err := api.NewPaginatedRequester(cliConnection, spaceRequestFactory)
	if err != nil {
		return err
	}

	spaces, err := resource_mapper.Spaces(
		models.SpacesParser{},
		spacePaginatedRequester,
	)
	if err != nil {
		return err
	}

	spaceMap := make(map[string]models.Space)
	for _, space := range spaces {
		spaceMap[space.Guid] = space
	}

	stopped, warnings, errors := exe.restartApps(cliConnection, apps, spaceMap)
	exe.RestartAppsUI.AfterAll(len(apps), stopped, warnings, errors)

	return nil
}

type restartAppFunc func(appPrinter *displayhelpers.AppPrinter, appRestarter AppRestarter) int

func (exe *RestartAppsExecutor) RestartApp(
	appPrinter *displayhelpers.AppPrinter,
	appRestarter AppRestarter,
) int {
	if appPrinter.App.State == models.Stopped {
		return Stopped
	}

	exe.RestartAppsUI.BeforeEach(appPrinter)

	waitTime := 1 * time.Minute
	timeout := os.Getenv("CF_STARTUP_TIMEOUT")
	if timeout != "" {
		t, err := strconv.Atoi(timeout)

		if err == nil {
			waitTime = time.Duration(float32(t)) * time.Second
		}
	}

	_, err := appRestarter.Restart(appPrinter.App.Guid)
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorized") {
			exe.RestartAppsUI.UserWarning(appPrinter)
			return Warning
		} else {
			exe.RestartAppsUI.FailRestart(appPrinter, err)
			return Err
		}
	}

	printDot := time.NewTicker(5 * time.Second)
	go func() {
		for range printDot.C {
			exe.RestartAppsUI.DuringEach(appPrinter)
		}
	}()

	time.Sleep(waitTime)

	printDot.Stop()

	exe.RestartAppsUI.CompletedEach(appPrinter)

	return Success
}

func (exe *RestartAppsExecutor) restartApps(cliConnection api.Connection, apps models.Applications, spaceMap map[string]models.Space) (int, int, int) {
	runningAppsChan := generateAppsChan(apps)
	outputsChan, waitDone := processAppsChan(cliConnection, spaceMap, exe.RestartApp, runningAppsChan, len(apps))

	waitDone.Wait()
	close(outputsChan)

	return outputAppsChan(outputsChan)
}

func generateAppsChan(apps models.Applications) chan models.Application {
	runningAppsChan := make(chan models.Application)
	go func() {
		defer close(runningAppsChan)
		for _, app := range apps {
			runningAppsChan <- app
		}
	}()

	return runningAppsChan
}

func processAppsChan(
	cliConnection api.Connection,
	spaceMap map[string]models.Space,
	restart restartAppFunc,
	appsChan chan models.Application,
	outputSize int) (chan int, *sync.WaitGroup) {
	var waitDone sync.WaitGroup

	output := make(chan int, outputSize)

	restarter := NewAppRestarter(cliConnection)

	waitDone.Add(1)

	go func() {
		defer waitDone.Done()

		for app := range appsChan {
			a := &displayhelpers.AppPrinter{
				App:    app,
				Spaces: spaceMap,
			}
			output <- restart(a, restarter)
		}
	}()

	return output, &waitDone
}

func outputAppsChan(outputsChan chan int) (int, int, int) {
	warnings := 0
	errors := 0
	stopped := 0

	for result := range outputsChan {
		switch result {
		case Warning:
			warnings++
		case Err:
			errors++
		case Stopped:
			stopped++
		default:
		}
	}
	return stopped, warnings, errors
}
