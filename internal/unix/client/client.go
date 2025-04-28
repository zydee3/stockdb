package client

import (
	"os"

	"github.com/urfave/cli"

	"github.com/zydee3/stockdb/internal/common/logger"
)

func Init() {
	app := cli.NewApp()
	app.Name = "stockctl"
	app.Usage = "Command-line tool for StockDB"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		applyYamlCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Errorf("Error: %s", err.Error())
		os.Exit(1)
	}
}
