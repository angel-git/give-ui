package logger

import (
	"log"
	"os"
)

const logFileName = "give-ui.error.log"

type ErrorFileLogger struct {
	logFile *os.File
}

func SetupLogger() *ErrorFileLogger {
	logFile, _ := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	log.SetOutput(logFile)
	return &ErrorFileLogger{
		logFile: logFile,
	}
}

func (fileLogger *ErrorFileLogger) Close() error {
	return fileLogger.logFile.Close()
}
