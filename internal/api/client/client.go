package apiclient

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
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
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
