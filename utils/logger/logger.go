package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/raphael-p/beango/config"
)

type MyLogger struct {
	stdOutLogger *log.Logger
	fileLogger   *log.Logger
	logLevel     logLevel
}

type logLevel uint8

const (
	logLevelTrace logLevel = iota
	logLevelDebug
	logLevelInfo
	logLevelWarning
	logLevelError
)

type nowFunc func() string

var now nowFunc = func() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000")
}

func newLogger(out io.Writer) *log.Logger {
	return log.New(out, "", 0)
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

func Init() {
	logger.logLevel = logLevel(config.Values.Logger.DefaultLevel.Value)
	logFile, err := openLogFile(config.Values.Logger.Directory, config.Values.Logger.Filename)
	if err != nil {
		panic(err.Error())
	}
	logger.fileLogger = newLogger(logFile)
}

func Close() {
	if logger.fileLogger == nil {
		return
	}
	if file, ok := logger.fileLogger.Writer().(*os.File); ok {
		Trace("closing log file")
		err := file.Close()
		if err != nil {
			Error(fmt.Sprint("failed to close log file: ", err))
		}
	}
}

func logMessage(level string, ansiColour string, message string) {
	reset := "\033[0m"
	time := now()
	logger.stdOutLogger.Printf("%s %s[%s]%s %s", time, ansiColour, level, reset, message)
	if logger.fileLogger != nil {
		logger.fileLogger.Printf("%s [%s] %s", time, level, message)
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
