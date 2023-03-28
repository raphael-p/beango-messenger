package utils

import (
	"log"
	"os"
)

var fl *fatalLogger

type fatalLogger struct {
	logger *log.Logger
}

func CreateFatalLogger() {
	fl = &fatalLogger{log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)}
}

func (fl *fatalLogger) Log(message string) {
	reset := "\033[0m"
	red := "\033[31;1m"
	fl.logger.Fatalf("%s[%s]%s %s", red, "FATAL_ERROR", reset, message)
}
