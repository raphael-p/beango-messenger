package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func mockConsoleLogger(t *testing.T) *bytes.Buffer {
	oldStdOutLogger := *logger.stdOutLogger
	t.Cleanup(func() { logger.stdOutLogger = &oldStdOutLogger })
	var buf bytes.Buffer
	logger.stdOutLogger = newLogger(&buf)
	return &buf
}

func mockFileLogger(t *testing.T) *bytes.Buffer {
	var oldFileLogger *log.Logger
	if logger.fileLogger != nil {
		loggerCopy := *logger.fileLogger
		oldFileLogger = &loggerCopy
	}
	t.Cleanup(func() { logger.fileLogger = oldFileLogger })
	var buf bytes.Buffer
	logger.fileLogger = newLogger(&buf)
	return &buf
}

func checkLog(t *testing.T, buf *bytes.Buffer, level, colour, message string, isConsole bool) {
	log := buf.String()
	var xLog string
	if isConsole {
		reset := "\033[0m"
		xLog = fmt.Sprintf("%s %s[%s]%s %s\n", now(), colour, level, reset, message)
	} else {
		xLog = fmt.Sprintf("%s [%s] %s\n", now(), level, message)
	}
	assert.Equals(t, log, xLog)
}

func TestMain(m *testing.M) {
	originalNow := now
	now = func() string { return "2020-06-20 18:52:13.303" }
	defer func() { now = originalNow }()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestLogMessage(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		consoleBuffer := mockConsoleLogger(t)
		fileBuffer := mockFileLogger(t)
		level := "TEST"
		message := "a test log"
		purple := "\033[0;35m"
		logMessage(level, purple, message)
		checkLog(t, consoleBuffer, level, purple, message, true)
		checkLog(t, fileBuffer, level, "", message, false)
	})
}

func TestLogFunctions(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		testCases := []struct {
			name, level, colour string
			logFunction         func(string)
		}{
			{"Trace", "TRACE", "", Trace},
			{"Debug", "DEBUG", "\033[34m", Debug},
			{"Info", "INFO", "\033[36m", Info},
			{"Warning", "WARNING", "\033[33;1m", Warning},
			{"Error", "ERROR", "\033[31;1m", Error},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buf := mockConsoleLogger(t)
				message := "a test log"
				testCase.logFunction(message)
				checkLog(t, buf, testCase.level, testCase.colour, message, true)
			})
		}
	})

	t.Run("LogLevelTooHigh", func(t *testing.T) {
		testCases := []struct {
			name        string
			logLevel    logLevel
			logFunction func(string)
		}{
			{"Trace", logLevelTrace, Trace},
			{"Debug", logLevelDebug, Debug},
			{"Info", logLevelInfo, Info},
			{"Warning", logLevelWarning, Warning},
			{"Error", logLevelError, Error},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				oldLogLevel := logger.logLevel
				t.Cleanup(func() { logger.logLevel = oldLogLevel })
				logger.logLevel = testCase.logLevel + 1
				buf := mockConsoleLogger(t)
				testCase.logFunction("a test log")
				assert.Equals(t, buf.String(), "")
			})
		}
	})
}
