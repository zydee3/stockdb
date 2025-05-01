package main

import (
	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/common/version"
	"github.com/zydee3/stockdb/internal/unix/client"
)

func main() {
	logger.Infof("Starting StockCtl. Version: %s", version.GetVersion())

	client.Init()
}
