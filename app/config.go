package app

import (
	"os"
	"path/filepath"

	"github.com/prgres/go2git-switch/config"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

/* --- cli config helper funcs --- */
func ConfigReload(c *cli.Context) error {
	return config.New(c.String("config"))
}

func ConfigSave(c *cli.Context) error {
	return config.Save()
}

func configInitFile(c *cli.Context) error {
	configPath := c.String("config")

	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		return err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		zap.S().Debugf("Config file not found. Creating: %s", configPath)
		viper.SetConfigType("yaml")
		if err := viper.WriteConfigAs(configPath); err != nil {
			return err
		}

		return nil
	}

	return nil
}
