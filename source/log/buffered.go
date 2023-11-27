package log

// A buffered logger provider. It writes messages to a target logger once the buffer is filled.
type BufferedLogger struct {
	Logger

	maxLevel         Level
	target           Logger
	messageCount     int
	bufferedMessages []messageData
}

// Creates a new buffered logger, targetting the specified logger.
func Buffered(target Logger, bufferSize int) *BufferedLogger {
	return &BufferedLogger{
		maxLevel:         ErrorLevel,
		target:           target,
		messageCount:     0,
		bufferedMessages: make([]messageData, bufferSize),
	}
}

func (logger *BufferedLogger) Initialize() error {
	return nil
}

func (logger *BufferedLogger) Finalize() {
	logger.Flush()
}

func (logger *BufferedLogger) SetMaxLevel(level Level) {
	logger.maxLevel = level
	logger.target.SetMaxLevel(level)
}

func (logger *BufferedLogger) Write(level Level, tag string, format string, values ...interface{}) {
	if level > logger.maxLevel {
		return
	}

	logger.bufferedMessages[logger.messageCount] = messageData{
		level:  level,
		tag:    tag,
		format: format,
		values: values,
	}

	logger.messageCount++

	if logger.messageCount == len(logger.bufferedMessages) {
		logger.Flush()
	}
}

// Writes the current buffered messages to the target logger.
func (logger *BufferedLogger) Flush() {
	for index := 0; index < logger.messageCount; index++ {
		logger.target.Write(
			logger.bufferedMessages[index].level,
			logger.bufferedMessages[index].tag,
			logger.bufferedMessages[index].format,
			logger.bufferedMessages[index].values...,
		)
	}

	logger.messageCount = 0
}
