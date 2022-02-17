package main

import (
	"fmt"
	"os"
	"time"

	"github.com/prgres/go2git-switch/app"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "go2git-switch",
		Usage:                "TBD",
		Compiled:             time.Now(),
		Flags:                app.Flags(),
		Authors: []*cli.Author{
			{
				Name: "M. Więcek",
			},
		},
		Action: app.ProfileSelect,
		Before: app.Before,
		After:  app.After,
		CommandNotFound: func(context *cli.Context, err string) {
			cli.ShowAppHelp(context)
		},
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "TBD",
				Before: app.ConfigReload,
				After:  app.CofnigSave,
				Action: app.ProfileAdd,
			},
			{

				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "TBD",
				Before:  app.ConfigReload,
				After:   app.CofnigSave,
				Action:  app.ProfileRemove,
			},
			{
				Name:   "current",
				Usage:  "TBD",
				Before: app.ConfigReload,
				Action: app.ProfileCurrent,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		zap.S().Fatal(fmt.Errorf("app failed: %w", err))
	}
}
