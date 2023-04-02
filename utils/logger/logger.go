package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/raphael-p/beango/config"
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

var logger *MyLogger = &MyLogger{newLogger(os.Stdout), nil, logLevel(config.Values.Logger.DefaultLevel)}

func newLogger(out io.Writer) *log.Logger {
	return log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func logMessage(level string, ansiColour string, message string) {
	reset := "\033[0m"
	logger.stdOutLogger.Printf("%s[%s]%s %s", ansiColour, level, reset, message)
	if logger.fileLogger != nil {
		logger.fileLogger.Printf("[%s] %s", level, message)
	}
}

func Trace(message string) {
	if logger.logLevel <= logLevelTrace {
		logMessage("TRACE", "", message)
	}
}

func Debug(message string) {
	if logger.logLevel <= logLevelDebug {
		logMessage("DEBUG", "\033[34m", message)
	}
}

func Info(message string) {
	if logger.logLevel <= logLevelInfo {
		logMessage("INFO", "\033[36m", message)
	}
}

func Warning(message string) {
	if logger.logLevel <= logLevelWarning {
		logMessage("WARNING", "\033[33;1m", message)
	}
}

func Error(message string) {
	if logger.logLevel <= logLevelError {
		logMessage("ERROR", "\033[31;1m", message)
	}
}

func Fatal(message string) {
	if logger.logLevel <= logLevelError {
		logMessage("FATAL ERROR", "\033[31;1m", message)
	}
	os.Exit(1)
}

func OpenLogFile(fail func(string)) {
	logDirectory := config.Values.Logger.Directory
	logFileName := config.Values.Logger.Filename
	path := filepath.Join(logDirectory, logFileName)
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = os.MkdirAll(logDirectory, 0755)
		if err != nil {
			fail(fmt.Sprint("failed to create log directory: ", err))
		}

		logFile, err = os.Create(path)
		if err != nil {
			fail(fmt.Sprint("failed to create log file: ", err))
		}
	}

	logger.fileLogger = newLogger(logFile)
}

func CloseLogFile() {
	if file, ok := logger.fileLogger.Writer().(*os.File); ok {
		file.Close()
	}
}
