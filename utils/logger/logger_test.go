package logger

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

// TODO: test each log level method
// TODO: test that messages dont get logged if log level is higher

func makeMessage(level, colour, message string) string {
	reset := "\033[0m"
	return fmt.Sprintf("%s %s[%s]%s %s\n", now(), colour, level, reset, message)
}

func TestMain(m *testing.M) {
	originalNow := now
	now = func() string { return "2020-06-20 18:52:13.303" }
	defer func() { now = originalNow }()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestLogMessage(t *testing.T) {
	oldStdOutLogger := *logger.stdOutLogger
	t.Cleanup(func() { logger.stdOutLogger = &oldStdOutLogger })
	var buf bytes.Buffer
	logger.stdOutLogger = newLogger(&buf)
	level := "TEST"
	message := "a test log"
	purple := "\033[0;35m"

	logMessage(level, purple, message)
	log := buf.String()
	xLog := makeMessage(level, purple, message)
	assert.Equals(t, log, xLog)
}
