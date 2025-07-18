package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qvcloud/go-project-template/internal/di"
	"github.com/qvcloud/gopkg/version"
	"github.com/urfave/cli/v3"
)

var (
	appName      = "xxName"
	appDesc      = ""
	appAuthor    = "yzimhao"
	appCopyright = "https://github.com/qvcloud/go-project-template"
)

func main() {
	_ = godotenv.Load()

	cmd := &cli.Command{
		Name:        appName,
		Description: appDesc,
		Usage:       "",
		Action: func(_ context.Context, cmd *cli.Command) error {
			app := di.App()
			app.Run()
			return nil
		},
		Commands: []*cli.Command{
			cmdVerison(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func cmdVerison() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Display version info.",
		Action: func(ctx context.Context, c *cli.Command) error {
			version.ShowVersion()
			return nil
		},
	}
}
