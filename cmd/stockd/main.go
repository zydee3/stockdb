package main

import (
	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/common/version"
	"github.com/zydee3/stockdb/internal/daemon"
)

func main() {
	logger.Infof("Starting StockD. Version: %s", version.GetVersion())

	daemon.Init()
}
