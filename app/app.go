package app

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultConfigPath = fmt.Sprintf("%s/.config/go2git-switch.yaml", os.Getenv("HOME"))
)

func Before(c *cli.Context) error {
	if err := handleFlags(c); err != nil {
		return err
	}

	if err := ConfigReload(c); err != nil {
		return err
	}

	return nil
}

func After(c *cli.Context) error {
	if c.String("config") != defaultConfigPath {
		prompt := promptui.Select{
			Label:    fmt.Sprintf(">>> Write current config to default path? (%s)", defaultConfigPath),
			HideHelp: true,
			Items:    []string{"yes", "no"},
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Name }}",
				Selected: fmt.Sprintf(`{{ "%s " }} {{ . | green }}`, promptui.IconGood),
			},
		}

		_, res, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("promp run failed %w", err)
		}

		if res == "yes" {
			if err := viper.WriteConfigAs(defaultConfigPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   defaultConfigPath,
			Usage:   "TBD",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   false,
			Usage:   "if print verbose logs",
		},
		&cli.StringFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Value:   "global",
			Usage:   "TBD",
		},
	}
}
func handleFlags(c *cli.Context) error {
	if c.Bool("verbose") {
		cfg := zap.Config{
			Encoding:    "console",
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			OutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:   "message",
				LevelKey:     "level",
				EncodeLevel:  zapcore.CapitalLevelEncoder,
				TimeKey:      "time",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				CallerKey:    "caller",
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}

		log, err := cfg.Build()
		if err != nil {
			return err
		}

		zap.ReplaceGlobals(log)
	}

	return nil
}
