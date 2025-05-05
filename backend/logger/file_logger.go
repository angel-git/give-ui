package logger

import (
	"fmt"
	"log"
	"os"
)

const logFileName = "give-ui.error.log"

type FileLogger struct {
	logFile *os.File
	logger  *log.Logger
}

func SetupLogger() *FileLogger {
	logFile, _ := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	log.SetOutput(logFile)
	return &FileLogger{
		logFile: logFile,
		logger:  log.New(logFile, "", log.LstdFlags),
	}
}

func (l *FileLogger) Print(msg string) {
	l.writeMsg("[INFO] " + msg)
}
func (l *FileLogger) Trace(msg string) {
	l.writeMsg("[TRACE] " + msg)
}
func (l *FileLogger) Debug(msg string) {
	l.writeMsg("[DEBUG] " + msg)
}
func (l *FileLogger) Info(msg string) {
	l.writeMsg("[INFO] " + msg)
}
func (l *FileLogger) Warning(msg string) {
	l.writeMsg("[WARN] " + msg)
}
func (l *FileLogger) Error(msg string) {
	l.writeMsg("[ERROR] " + msg)
}
func (l *FileLogger) Fatal(msg string) {
	l.writeMsg("[FATAL] " + msg)
}
func (l *FileLogger) writeMsg(msg string) {
	if err := l.logger.Output(2, msg); err != nil {
		fmt.Println("Logger Error:", err) //nolint:forbidigo
	}
}

func (fileLogger *FileLogger) Close() error {
	return fileLogger.logFile.Close()
}
