package client

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/common/version"
)

func Init() {
	cmd := &cli.Command {
		Name: "stockctl",
		Description: "Command-line tool for StockDB",
		Version: version.GetVersion(),
		Commands: []*cli.Command{
			&applyYamlCommand,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Error("%w", err)
		os.Exit(1)
	}
}
