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

func newLogger(out io.Writer) *log.Logger {
	return log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

var logger *MyLogger = &MyLogger{newLogger(os.Stdout), nil, logLevel(0)}

func openLogFile(directory, name string) (*os.File, error) {
	path := filepath.Join(directory, name)
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory: %s", err)
		}

		logFile, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %s", err)
		}
	}
	return logFile, nil
}

func Init(fail func(string)) {
	logger.logLevel = logLevel(config.Values.Logger.DefaultLevel)
	logFile, err := openLogFile(config.Values.Logger.Directory, config.Values.Logger.Filename)
	if err != nil {
		fail(err.Error())
		return
	}
	logger.fileLogger = newLogger(logFile)
}

func Close() {
	if file, ok := logger.fileLogger.Writer().(*os.File); ok {
		file.Close()
	}
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
