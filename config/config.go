package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func New(configPath string) error {
	zap.L().Debug("Initializing config.")
	configFilename, path := parsePath(configPath)

	viper.SetConfigName(configFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Debug(fmt.Sprintf("Config file changed: %s", e.Name))
		read()
	})

	viper.WatchConfig()

	return viper.ReadInConfig()
}

func Reload(c *cli.Context) error {
	return New(c.String("config"))
}

/* --- helper func --- */
func parsePath(configPath string) (string, string) {
	pathParts := strings.Split(configPath, "/")
	configFilename := pathParts[len(pathParts)-1]
	path := func() string {
		og := strings.Join(pathParts[:len(pathParts)-1], "/")
		if og == "" {
			og = "."
		}
		return og
	}()

	return configFilename, path
}

func read() {
	zap.L().Debug("Loading profiles from config file.")
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatal(fmt.Errorf("reload config failed: %w", err))
	}
}

func Save() error {
	return viper.WriteConfig()
}
