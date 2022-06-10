package log

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var loggerInitialized = false

func init() {
	switch mode := os.Getenv("MODE"); mode {
	case "PROD":
		newLogger, err := zap.NewProduction()
		postConfigureLogger(newLogger, err)
	case "DEV":
		fallthrough
	default:
		newLogger, err := zap.NewDevelopment()
		postConfigureLogger(newLogger, err)
	}
}

func postConfigureLogger(logger *zap.Logger, err error) {
	if err != nil {
		log.Println("Logger couldn't be initialized", err)
	} else {
		loggerInitialized = true
		logger = logger.WithOptions(zap.AddCallerSkip(1)).WithOptions(zap.AddStacktrace(zap.FatalLevel))
		zap.ReplaceGlobals(logger)
	}
}

func Debug(v ...interface{}) {
	if loggerInitialized {
		logger := zap.S()
		defer logger.Sync()
		logger.Debug(v...)
	} else {
		log.Println(v...)
	}
}

func Error(v ...interface{}) {
	if loggerInitialized {
		logger := zap.S()
		defer logger.Sync()
		logger.Error(v...)
	} else {
		log.Println(v...)
	}
}
