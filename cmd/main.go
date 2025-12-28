package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qvcloud/go-project-template/internal/di"
	"github.com/qvcloud/go-project-template/internal/di/provider"
	"github.com/qvcloud/gopkg/version"
	"github.com/urfave/cli/v3"
)

var (
	appName = "go-project-template"
	appDesc = "A starter template for Go web services"
)

func main() {
	_ = godotenv.Load() //nolint:errcheck

	cmd := &cli.Command{
		Name:        appName,
		Description: appDesc,
		Usage:       "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				Value:   "",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			configPath := cmd.String("config")
			v := provider.NewConfig(configPath)
			app := di.App(v)
			app.Run()
			return nil
		},
		Commands: []*cli.Command{
			cmdVersion(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func cmdVersion() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Display version info.",
		Action: func(_ context.Context, _ *cli.Command) error {
			version.ShowVersion()
			return nil
		},
	}
}
