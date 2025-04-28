package main

import (
	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/daemon"
)

//nolint:gochecknoglobals // gochecknoglobals
var (
	version   = ""
	gitCommit = ""
)

func main() {
	logger.SetupLogger()
	logger.Infof("Starting StockDB. Version: %s. Commit: %s.", version, gitCommit)

	daemon.Init()
}
