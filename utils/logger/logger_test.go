package logger

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func makeMessage(level, colour, message string) string {
	reset := "\033[0m"
	return fmt.Sprintf("%s %s[%s]%s %s\n", now(), colour, level, reset, message)
}

func getLogBuffer(t *testing.T) *bytes.Buffer {
	oldStdOutLogger := *logger.stdOutLogger
	t.Cleanup(func() { logger.stdOutLogger = &oldStdOutLogger })
	var buf bytes.Buffer
	logger.stdOutLogger = newLogger(&buf)
	return &buf
}

func checkLog(t *testing.T, buf *bytes.Buffer, level, colour, message string) {
	log := buf.String()
	xLog := makeMessage(level, colour, message)
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
		buf := getLogBuffer(t)
		level := "TEST"
		message := "a test log"
		purple := "\033[0;35m"
		logMessage(level, purple, message)
		checkLog(t, buf, level, purple, message)
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
				buf := getLogBuffer(t)
				message := "a test log"
				testCase.logFunction(message)
				checkLog(t, buf, testCase.level, testCase.colour, message)
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
				buf := getLogBuffer(t)
				testCase.logFunction("a test log")
				assert.Equals(t, buf.String(), "")
			})
		}
	})
}
