package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/raphael-p/beango/config"
)

const MAX_FILE_COUNT = 1000
const MAX_FILE_BYTES = 10 * 1024 * 1024
const MAX_MESSAGE_BYTES = 10 * 1024

type MyLogger struct {
	StdOutLogger *log.Logger
	FileLogger   *log.Logger
	logLevel     logLevel
	cumBytes     int64
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

func NewLogger(out io.Writer) *log.Logger {
	return log.New(out, "", 0)
}

var Logger *MyLogger = &MyLogger{NewLogger(os.Stdout), nil, logLevel(0), 0}

func openLogFile() *os.File {
	directory := config.Values.Logger.Directory
	name := config.Values.Logger.Filename
	path := filepath.Join(directory, generateFilename(name, directory))
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			panic(fmt.Sprintf("failed to create log directory: %s", err))
		}

		logFile, err = os.Create(path)
		if err != nil {
			panic(fmt.Sprintf("failed to create log file: %s", err))
		}
	}
	return logFile
}

func generateFilename(name, directory string) string {
	filenameArr := strings.SplitN(name, ".", 2)
	basename := filenameArr[0]
	extension := ".log"
	if len(filenameArr) > 1 && filenameArr[1] != "" {
		extension = "." + filenameArr[1]
	}
	basename += "-" + time.Now().UTC().Format("20060102")

	count := 1
	for {
		filename := fmt.Sprint(basename, "-", count, extension)
		path := filepath.Join(directory, filename)
		_, err := os.Stat(path)
		if err != nil {
			return filename
		}
		count++
		if count > MAX_FILE_COUNT {
			panic(fmt.Sprint(
				"reached maximum allowed number of log files: ",
				MAX_FILE_COUNT,
			))
		}
	}
}

func Init() {
	Logger.logLevel = logLevel(config.Values.Logger.DefaultLevel.Value)
	Logger.FileLogger = NewLogger(openLogFile())
	Trace("file logger initialised")
}

func Close() {
	if Logger.FileLogger == nil {
		return
	}
	if file, ok := Logger.FileLogger.Writer().(*os.File); ok {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprint("failed to close log file: ", err))
		}
		Trace("log file closed")
	}
}

func logMessage(level string, ansiColour string, message string) {
	reset := "\033[0m"
	time := now()
	if len(message) > MAX_MESSAGE_BYTES {
		message = message[:MAX_MESSAGE_BYTES]
	}
	Logger.StdOutLogger.Printf("%s %s[%s]%s %s", time, ansiColour, level, reset, message)
	if Logger.FileLogger != nil {
		if Logger.cumBytes += int64(len(message)); Logger.cumBytes > MAX_FILE_BYTES {
			rollover()
		}
		Logger.FileLogger.Printf("%s [%s] %s", time, level, message)
	}
}

func rollover() {
	if file, ok := Logger.FileLogger.Writer().(*os.File); ok {
		if stats, err := file.Stat(); err == nil && stats.Size() > MAX_FILE_BYTES {
			Logger.cumBytes = 0
			Trace(fmt.Sprintf(
				"reached maximum log file size (%d bytes), rolling over",
				MAX_FILE_BYTES,
			)) // can cause infinite loop if MAX_FILE_BYTES is too low
			Close()
			Logger.FileLogger = NewLogger(openLogFile())
			Trace("new log file opened")
		}
	}
}

func Trace(message string) {
	if Logger.logLevel <= logLevelTrace {
		logMessage("TRACE", "", message)
	}
}

func Debug(message string) {
	if Logger.logLevel <= logLevelDebug {
		logMessage("DEBUG", "\033[34m", message)
	}
}

func Info(message string) {
	if Logger.logLevel <= logLevelInfo {
		logMessage("INFO", "\033[36m", message)
	}
}

func Warning(message string) {
	if Logger.logLevel <= logLevelWarning {
		logMessage("WARNING", "\033[33;1m", message)
	}
}

func Error(message string) {
	if Logger.logLevel <= logLevelError {
		buf := make([]byte, 1<<16)
		n := runtime.Stack(buf, false)
		stackTrace := strings.ReplaceAll(string(buf[:n-1]), "\n", "\n\t")
		message += fmt.Sprintf("\n\tstack trace: %s", stackTrace)
		logMessage("ERROR", "\033[31;1m", message)
	}
}
