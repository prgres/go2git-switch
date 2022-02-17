package app

import (
	"github.com/prgres/go2git-switch/config"
	"github.com/urfave/cli/v2"
)

/* --- cli config helper funcs --- */
func ConfigReload(c *cli.Context) error {
	return config.New(c.String("config"))
}

func CofnigSave(c *cli.Context) error {
	return config.Save()
}
