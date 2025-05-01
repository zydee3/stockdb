package logger

import (
	"fmt"
	"log/slog"
	"os"
)

//nolint:gochecknoglobals // gochecknoglobals
var logger = slog.Default()

func SetupLogger() {
	// TODO: Oscar - Load logging info from config file and cli options
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func Debug(msg string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Debug(msg, args...)
}

func Debugf(format string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Debug(fmt.Sprintf(format, args...))
}

func Info(msg string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Info(msg, args...)
}

func Infof(format string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Info(fmt.Sprintf(format, args...))
}

func Error(msg string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Error(msg, args...)
}

func Errorf(format string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Error(fmt.Sprintf(format, args...))
}

func Warn(msg string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Warn(msg, args...)
}

func Warnf(format string, args ...any) {
	//nolint:sloglint // sloglint
	logger.Warn(fmt.Sprintf(format, args...))
}
