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
	buf := getLogBuffer(t)
	level := "TEST"
	message := "a test log"
	purple := "\033[0;35m"
	logMessage(level, purple, message)
	checkLog(t, buf, level, purple, message)
}

func logTest(t *testing.T, level, colour string, fn func(string)) {
	buf := getLogBuffer(t)
	message := "a test log"
	fn(message)
	checkLog(t, buf, level, colour, message)
}
func TestTrace(t *testing.T)   { logTest(t, "TRACE", "", Trace) }
func TestDebug(t *testing.T)   { logTest(t, "DEBUG", "\033[34m", Debug) }
func TestInfo(t *testing.T)    { logTest(t, "INFO", "\033[36m", Info) }
func TestWarning(t *testing.T) { logTest(t, "WARNING", "\033[33;1m", Warning) }
func TestError(t *testing.T)   { logTest(t, "ERROR", "\033[31;1m", Error) }

func logTestBelowLevel(t *testing.T, level logLevel, fn func(string)) {
	oldLogLevel := logger.logLevel
	t.Cleanup(func() { logger.logLevel = oldLogLevel })
	logger.logLevel = level + 1
	buf := getLogBuffer(t)
	fn("a test log")
	assert.Equals(t, buf.String(), "")
}
func TestTrace_LogLevelIsHigher(t *testing.T)   { logTestBelowLevel(t, logLevelTrace, Trace) }
func TestDebug_LogLevelIsHigher(t *testing.T)   { logTestBelowLevel(t, logLevelDebug, Debug) }
func TestInfo_LogLevelIsHigher(t *testing.T)    { logTestBelowLevel(t, logLevelInfo, Info) }
func TestWarning_LogLevelIsHigher(t *testing.T) { logTestBelowLevel(t, logLevelWarning, Warning) }
func TestError_LogLevelIsHigher(t *testing.T)   { logTestBelowLevel(t, logLevelError, Error) }
