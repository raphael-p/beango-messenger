package logger

import (
	"bytes"
	"testing"
)

func MockConsoleLogger(t *testing.T) *bytes.Buffer {
	oldWriter := Logger.StdOutLogger.Writer()
	t.Cleanup(func() { Logger.StdOutLogger = NewLogger(oldWriter) })
	var buf bytes.Buffer
	Logger.StdOutLogger = NewLogger(&buf)
	return &buf
}

func MockFileLogger(t *testing.T) *bytes.Buffer {
	if Logger.FileLogger != nil {
		oldWriter := Logger.FileLogger.Writer()
		t.Cleanup(func() { Logger.FileLogger = NewLogger(oldWriter) })
	} else {
		t.Cleanup(func() { Logger.FileLogger = nil })
	}
	var buf bytes.Buffer
	Logger.FileLogger = NewLogger(&buf)
	return &buf
}
