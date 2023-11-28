package log

type BroadcastLogger struct {
	Logger

	maxLevel Level
	targets  []Logger
}

func Broadcast(targets ...Logger) *BroadcastLogger {
	var logger = BroadcastLogger{
		maxLevel: ErrorLevel,
		targets:  make([]Logger, 0),
	}

	if targets != nil {
		logger.targets = append(logger.targets, targets...)
	}

	return &logger
}

func (logger *BroadcastLogger) Initialize() error {
	return nil
}

func (logger *BroadcastLogger) Finalize() {
	// Empty
}

func (logger *BroadcastLogger) SetMaxLevel(level Level) {
	logger.maxLevel = level

	for _, target := range logger.targets {
		target.SetMaxLevel(level)
	}
}

func (logger *BroadcastLogger) Attach(target Logger) {
	if target == nil {
		return
	}

	logger.targets = append(logger.targets, target)
}

func (logger *BroadcastLogger) Write(level Level, tag string, format string, values ...interface{}) {
	if level > logger.maxLevel {
		return
	}

	for _, target := range logger.targets {
		target.Write(level, tag, format, values...)
	}
}
