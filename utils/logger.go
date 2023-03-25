package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type MyLogger struct {
	stdOutLogger *log.Logger
	fileLogger   *log.Logger
}

func NewLogger(logDirectory string, logFileName string) (*MyLogger, error) {
	path := filepath.Join(logDirectory, logFileName)
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = os.MkdirAll(logDirectory, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory: %s", err)
		}

		logFile, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %s", err)
		}
	}
	return &MyLogger{
		log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds/1000),
		log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds/1000),
	}, nil
}

func (l *MyLogger) Close() {
	if file, ok := l.fileLogger.Writer().(*os.File); ok {
		file.Close()
	}
}

func (l *MyLogger) log(level string, ansiColour string, message string) {
	reset := "\033[0m"
	l.stdOutLogger.Printf("%s[%s]%s %s", ansiColour, level, reset, message)
	l.fileLogger.Printf("[%s] %s", level, message)
}

func (l *MyLogger) Trace(message string) {
	l.log("TRACE", "", message)
}

func (l *MyLogger) Debug(message string) {
	l.log("DEBUG", "\033[34m", message)
}

func (l *MyLogger) Info(message string) {
	l.log("INFO", "\033[36m", message)
}

func (l *MyLogger) Warning(message string) {
	l.log("WARNING", "\033[33;1m", message)
}

func (l *MyLogger) Error(message string) {
	l.log("ERROR", "\033[31;1m", message)
}
