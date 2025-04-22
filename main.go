package main

import (
	"github.com/zydee3/stockdb/internal/common/logging"
)

var (
	version = ""
	gitCommit = ""
)

func main() {
	logger.SetupLogger()
	logger.Infof("Starting StockDB. Version: %s. Commit: %s.", version, gitCommit)
}
