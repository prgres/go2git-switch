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
		Usage:                "little helper to easily between git profiles",
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
				Usage:  "add a new profile to configfile",
				Before: app.ConfigReload,
				After:  app.ConfigSave,
				Action: app.ProfileAdd,
			},
			{
				Name:    "edit",
				Usage:   "edit existing profile",
				Aliases: []string{"e"},
				Before:  app.ConfigReload,
				After:   app.ConfigSave,
				Action:  app.ProfileEdit,
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "list all profiles and select one to remove from configfile",
				Before:  app.ConfigReload,
				After:   app.ConfigSave,
				Action:  app.ProfileRemove,
			},
			{
				Name:   "current",
				Usage:  "show current git profile",
				Before: app.ConfigReload,
				Action: app.ProfileCurrent,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		zap.S().Fatal(fmt.Errorf("app failed: %w", err))
	}
}
