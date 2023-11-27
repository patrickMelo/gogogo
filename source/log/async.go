package log

import "sync"

// An async logger provider. It writes non-blocking messages to a target logger.
type AsyncLogger struct {
	Logger

	maxLevel     Level
	target       Logger
	waitGroup    sync.WaitGroup
	writeChannel chan messageData
}

// Creates a new async logger, targetting the specified logger.
func Async(target Logger, bufferSize int) *AsyncLogger {
	return &AsyncLogger{
		maxLevel:     ErrorLevel,
		target:       target,
		writeChannel: make(chan messageData, bufferSize),
	}
}

func (logger *AsyncLogger) Initialize() error {
	go logger.asyncWrite()
	return nil
}

func (logger *AsyncLogger) Finalize() {
	logger.waitGroup.Wait()
	close(logger.writeChannel)
}

func (logger *AsyncLogger) SetMaxLevel(level Level) {
	logger.maxLevel = level
	logger.target.SetMaxLevel(level)
}

func (logger *AsyncLogger) Write(level Level, tag string, format string, values ...interface{}) {
	if level > logger.maxLevel || logger.writeChannel == nil {
		return
	}

	logger.waitGroup.Add(1)
	logger.writeChannel <- messageData{
		level:  level,
		tag:    tag,
		format: format,
		values: values,
	}
}

func (logger *AsyncLogger) asyncWrite() {
	for logger.writeChannel != nil {
		var message = <-logger.writeChannel
		logger.target.Write(message.level, message.tag, message.format, message.values...)
		logger.waitGroup.Done()
	}
}
