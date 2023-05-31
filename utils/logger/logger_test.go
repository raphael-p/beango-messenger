package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func checkLog(t *testing.T, buf *bytes.Buffer, level, colour, message string, isConsole bool) {
	log := buf.String()
	var xLog string
	if isConsole {
		reset := "\033[0m"
		xLog = fmt.Sprintf("%s %s[%s]%s %s\n", now(), colour, level, reset, message)
	} else {
		xLog = fmt.Sprintf("%s [%s] %s\n", now(), level, message)
	}
	if level != "ERROR" {
		assert.Equals(t, log, xLog)
	} else {
		assert.Equals(t, strings.Split(log, "\n")[0]+"\n", xLog)
		assert.Contains(t, log, "stack trace: ")
	}
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
		consoleBuffer := MockConsoleLogger(t)
		fileBuffer := MockFileLogger(t)
		level := "TEST"
		message := "a test log"
		purple := "\033[0;35m"
		logMessage(level, purple, message)
		checkLog(t, consoleBuffer, level, purple, message, true)
		checkLog(t, fileBuffer, level, "", message, false)
	})

	t.Run("Wrappers", func(t *testing.T) {
		testCases := []struct {
			name, level, colour string
			logLevel            logLevel
			logMessageWrapper   func(string)
		}{
			{"Trace", "TRACE", "", logLevelTrace, Trace},
			{"Debug", "DEBUG", "\033[34m", logLevelDebug, Debug},
			{"Info", "INFO", "\033[36m", logLevelInfo, Info},
			{"Warning", "WARNING", "\033[33;1m", logLevelWarning, Warning},
			{"Error", "ERROR", "\033[31;1m", logLevelError, Error},
		}

		t.Run("Normal", func(t *testing.T) {
			for _, testCase := range testCases {
				t.Run(testCase.name, func(t *testing.T) {
					buf := MockConsoleLogger(t)
					message := "a test log"
					testCase.logMessageWrapper(message)
					checkLog(t, buf, testCase.level, testCase.colour, message, true)
				})
			}
		})

		t.Run("LogLevelTooHigh", func(t *testing.T) {
			for _, testCase := range testCases {
				t.Run(testCase.name, func(t *testing.T) {
					oldLogLevel := Logger.logLevel
					t.Cleanup(func() { Logger.logLevel = oldLogLevel })
					Logger.logLevel = testCase.logLevel + 1
					buf := MockConsoleLogger(t)
					testCase.logMessageWrapper("a test log")
					assert.Equals(t, buf.String(), "")
				})
			}
		})
	})
}
