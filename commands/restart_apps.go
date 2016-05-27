package commands

import (
	"github.com/cloudfoundry-incubator/app-restarter/commands/errorhelpers"
	"github.com/cloudfoundry-incubator/app-restarter/resource_mapper"
	"github.com/cloudfoundry-incubator/app-restarter/ui"
)

type RestartAppsCommand struct {
	Organization string `short:"o" value-name:"ORG" description:"Organization to restrict the app restarts"`
	Space        string `short:"s" value-name:"SPACE" description:"Space in the targeted organization to restrict the app restarts"`
}

func (command RestartAppsCommand) Execute(flags []string) error {
	cliConnection := Context.CLIConnection

	err := errorhelpers.ErrorIfOrgAndSpacesSet(command.Organization, command.Space)
	if err != nil {
		return err
	}

	appsGetter, err := resource_mapper.NewAppsGetterFunc(cliConnection, command.Organization, command.Space)
	if err != nil {
		return err
	}

	restartAppsUI, err := ui.NewRestartApps(cliConnection, command.Organization, command.Space)
	if err != nil {
		return err
	}

	cmd := RestartAppsExecutor{
		AppsGetterFunc: appsGetter,
		RestartAppsUI:  &restartAppsUI,
	}

	return cmd.Execute(cliConnection)
}
