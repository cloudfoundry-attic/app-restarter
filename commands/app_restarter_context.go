package commands

import "github.com/cloudfoundry-incubator/app-restarter/api"

type AppRestarterContext struct {
	CLIConnection api.Connection

	RestartApps RestartAppsCommand `command:"restart-apps" description:"Restart all apps"`
}

var Context AppRestarterContext
