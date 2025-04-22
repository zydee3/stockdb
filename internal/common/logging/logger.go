package logger

import (
	"fmt"
        "log/slog"
	"os"
)

var logger *slog.Logger

func SetupLogger() {
        // TODO: Load logging info form config file and cli options
        logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Debugf(format string, args ...any) {
	logger.Debug(fmt.Sprintf(format, args...))
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Infof(format string, args ...any) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Error(msg string, args ...any) {
	logger.Error(msg, args)
}

func Errorf(format string, args ...any) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Warn(msg string, args ...any) {
        logger.Warn(msg, args...)
}

func Warnf(format string, args ...any) {
        logger.Warn(fmt.Sprintf(format, args...))
}
