package main

import (
	"github.com/zydee3/stockdb/internal/api/client"
	"github.com/zydee3/stockdb/internal/common/logger"
)

var (
	version   = ""
	gitCommit = ""
)

func main() {
	logger.SetupLogger()
	logger.Infof("Starting StockDB. Version: %s. Commit: %s.", version, gitCommit)

	client.Init()
}
