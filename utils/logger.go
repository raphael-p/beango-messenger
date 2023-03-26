package utils

import (
	"log"
	"os"
	"path/filepath"
)

type MyLogger struct {
	stdOutLogger *log.Logger
	fileLogger   *log.Logger
	logLevel     logLevel
}

type logLevel int

const (
	logLevelTrace logLevel = iota
	logLevelDebug
	logLevelInfo
	logLevelWarning
	logLevelError
)

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
	if l.logLevel <= logLevelTrace {
		l.log("TRACE", "", message)
	}
}

func (l *MyLogger) Debug(message string) {
	if l.logLevel <= logLevelDebug {
		l.log("DEBUG", "\033[34m", message)
	}
}

func (l *MyLogger) Info(message string) {
	if l.logLevel <= logLevelInfo {
		l.log("INFO", "\033[36m", message)
	}
}

func (l *MyLogger) Warning(message string) {
	if l.logLevel <= logLevelWarning {
		l.log("WARNING", "\033[33;1m", message)
	}
}

func (l *MyLogger) Error(message string) {
	if l.logLevel <= logLevelError {
		l.log("ERROR", "\033[31;1m", message)
	}
}

var Logger *MyLogger

func init() {
	logDirectory := "logs"          // TODO: make config variable
	logFileName := "server.log"     // TODO: make config variable
	defaultLogLevel := logLevelInfo // TODO: make config variable
	path := filepath.Join(logDirectory, logFileName)
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = os.MkdirAll(logDirectory, 0755)
		if err != nil {
			log.Fatalf("FATAL ERROR failed to create log directory: %s\n", err)
		}

		logFile, err = os.Create(path)
		if err != nil {
			log.Fatalf("FATAL ERROR failed to create log file: %s\n", err)
		}
	}

	Logger = &MyLogger{
		log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		defaultLogLevel,
	}
	Logger.Trace("logger created")
}
