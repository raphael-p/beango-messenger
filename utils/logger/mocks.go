package logger

import (
	"bytes"
	"log"
	"testing"
)

func MockConsoleLogger(t *testing.T) *bytes.Buffer {
	oldStdOutLogger := *Logger.StdOutLogger
	t.Cleanup(func() { Logger.StdOutLogger = &oldStdOutLogger })
	var buf bytes.Buffer
	Logger.StdOutLogger = NewLogger(&buf)
	return &buf
}

func MockFileLogger(t *testing.T) *bytes.Buffer {
	var oldFileLogger *log.Logger
	if Logger.FileLogger != nil {
		loggerCopy := *Logger.FileLogger
		oldFileLogger = &loggerCopy
	}
	t.Cleanup(func() { Logger.FileLogger = oldFileLogger })
	var buf bytes.Buffer
	Logger.FileLogger = NewLogger(&buf)
	return &buf
}
