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

// @title Go Project Template API
// @version 1.0
// @description This is a sample server for a Go project template.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
