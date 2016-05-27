package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/app-restarter/commands"
	"github.com/cloudfoundry-incubator/app-restarter/ui"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jessevdk/go-flags"

)

type AppRestarter struct{}

func (c *AppRestarter) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "app-restarter",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "restart-apps",
				HelpText: "Restart all apps",
				UsageDetails: plugin.Usage{
					Usage: `cf restart-apps [-o ORG | -s SPACE]

OPTIONS:
   -o      Organization to restrict the app restarts
   -s      Space in the targeted organization to restrict the app restarts`,
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(AppRestarter))
}

func (c *AppRestarter) Run(cliConnection plugin.CliConnection, args []string) {
	commands.Context.CLIConnection = cliConnection

	parser := flags.NewParser(&commands.AppRestarterContext{}, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	_, err := parser.ParseArgs(args)
	if err != nil {
		ui.SayFailed()
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
