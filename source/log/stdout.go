package log

import (
	"fmt"
	"time"
)

type StdoutLogger struct {
	Logger

	maxLevel Level
}

func Stdout() *StdoutLogger {
	return &StdoutLogger{
		maxLevel: ErrorLevel,
	}
}

func (logger *StdoutLogger) Initialize() error {
	return nil
}

func (logger *StdoutLogger) Finalize() {
	// Empty
}

func (logger *StdoutLogger) SetMaxLevel(level Level) {
	logger.maxLevel = level
}

func (logger *StdoutLogger) Write(level Level, tag string, format string, values ...interface{}) {
	if level > logger.maxLevel {
		return
	}

	fmt.Printf("[%s] (%c) [%s] %s\n", time.Now().Format("2006-01-02 15:04:05.000000"), level.LogID(), tag, fmt.Sprintf(format, values...))
}
