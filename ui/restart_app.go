package ui

import (
	"fmt"

	"github.com/cloudfoundry-incubator/app-restarter/api"
	"github.com/cloudfoundry/cli/cf/terminal"
)

type RestartApps struct {
	Username     string
	Organization string
	Space        string
}

func NewRestartApps(cliConnection api.Connection, organizationName string, spaceName string) (RestartApps, error) {
	username, err := cliConnection.Username()
	if err != nil {
		return RestartApps{}, err
	}

	if spaceName != "" {
		space, err := cliConnection.GetSpace(spaceName)
		if err != nil || space.Guid == "" {
			return RestartApps{}, err
		}
		organizationName = space.Organization.Name
	}

	return RestartApps{
		Username:     username,
		Organization: organizationName,
		Space:        spaceName,
	}, nil
}

func (c *RestartApps) BeforeAll() {
	switch {
	case c.Organization != "" && c.Space != "":
		fmt.Printf(
			"Restarting apps in org %s / %s as %s...\n",
			terminal.EntityNameColor(c.Organization),
			terminal.EntityNameColor(c.Space),
			terminal.EntityNameColor(c.Username),
		)
	case c.Organization != "":
		fmt.Printf(
			"Restarting apps in org %s as %s...\n",
			terminal.EntityNameColor(c.Organization),
			terminal.EntityNameColor(c.Username),
		)
	default:
		fmt.Printf(
			"Restarting apps as %s...\n",
			terminal.EntityNameColor(c.Username),
		)
	}
}

func (c *RestartApps) BeforeEach(app ApplicationPrinter) {
	fmt.Println()
	fmt.Printf(
		"Restarting app %s in org %s / space %s as %s...\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(c.Username),
	)
}

func (c *RestartApps) CompletedEach(app ApplicationPrinter) {
	fmt.Println()
	fmt.Printf(
		"Completed restarting app %s in org %s / space %s as %s\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(c.Username),
	)
}

func (c *RestartApps) DuringEach(app ApplicationPrinter) {
	fmt.Print(".")
}

func (c *RestartApps) AfterAll(attempts, warnings int, errors int) {
	successes := attempts - warnings - errors
	fmt.Println()
	fmt.Printf("Restarting completed: %d apps, %d errors, %d warnings\n", successes, errors, warnings)
}

func (c *RestartApps) UserWarning(app ApplicationPrinter) {
	fmt.Printf(
		"WARNING: No authorization to restart app %s in space %s / org %s as %s\n",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(c.Username),
	)
}

func (c *RestartApps) FailRestart(app ApplicationPrinter, err error) {
	fmt.Printf(
		"Error: Failed to restart app %s in space %s / org %s as %s: %s",
		terminal.EntityNameColor(app.Name()),
		terminal.EntityNameColor(app.Space()),
		terminal.EntityNameColor(app.Organization()),
		terminal.EntityNameColor(c.Username),
		terminal.EntityNameColor(err.Error()),
	)
}
