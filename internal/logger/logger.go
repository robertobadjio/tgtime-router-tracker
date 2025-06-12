package logger

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

var globalLogger log.Logger

func init() {
	globalLogger = log.NewLogfmtLogger(os.Stdout)
	globalLogger = log.With(globalLogger, "ts", log.DefaultTimestampUTC)
	globalLogger = level.NewFilter(globalLogger, level.AllowInfo())
}

// Log ...
func Log(keyVals ...interface{}) {
	_ = globalLogger.Log(keyVals...)
}

// Debug ...
func Debug(keyVals ...interface{}) {
	_ = level.Debug(globalLogger).Log(keyVals...)
}

// Info ...
func Info(keyVals ...interface{}) {
	_ = level.Info(globalLogger).Log(keyVals...)
}

// Warn ...
func Warn(keyVals ...interface{}) {
	_ = level.Warn(globalLogger).Log(keyVals...)
}

// Error ...
func Error(keyVals ...interface{}) {
	_ = level.Error(globalLogger).Log(keyVals...)
}

// Fatal ...
func Fatal(keyVals ...interface{}) {
	_ = level.Error(globalLogger).Log(keyVals...)
	os.Exit(1)
}
