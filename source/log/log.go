package log

// Defines the log level type.
type Level int

const (
	FatalLevel Level = iota
	ErrorLevel
	WarningLevel
	InformationLevel
	VerboseLevel
)

func (level Level) String() string {
	switch level {
	case FatalLevel:
		return "Fatal"
	case ErrorLevel:
		return "Error"
	case WarningLevel:
		return "Warning"
	case InformationLevel:
		return "Information"
	case VerboseLevel:
		return "Verbose"
	default:
		return "?"
	}
}

func (level Level) LogID() rune {
	switch level {
	case FatalLevel:
		return 'F'
	case ErrorLevel:
		return 'E'
	case WarningLevel:
		return 'W'
	case InformationLevel:
		return 'I'
	case VerboseLevel:
		return 'V'
	default:
		return '?'
	}
}

// Defines the interface for logging providers.
type Logger interface {
	Initialize() error
	Finalize()
	SetMaxLevel(level Level)
	Write(level Level, tag string, format string, values ...interface{})
}

// Initializes the global log infrastructure using the provided default logger.
func Initialize(defaultLogger Logger) error {
	_defaultLogger = defaultLogger
	return _defaultLogger.Initialize()
}

// Finalizes the default logger and the global log infrastructure.
func Finalize() {
	_defaultLogger.Finalize()
	_defaultLogger = nil
}

// Sets the default logger max level.
func SetMaxLevel(level Level) {
	if _defaultLogger == nil {
		return
	}

	_defaultLogger.SetMaxLevel(level)
}

// Logs a verbose message to the default logger.
func Verbose(tag string, format string, values ...interface{}) {
	if _defaultLogger == nil {
		return
	}

	_defaultLogger.Write(VerboseLevel, tag, format, values...)
}

// Logs a information message to the default logger.
func Information(tag string, format string, values ...interface{}) {
	if _defaultLogger == nil {
		return
	}

	_defaultLogger.Write(InformationLevel, tag, format, values...)
}

// Logs a warning message to the default logger.
func Warning(tag string, format string, values ...interface{}) {
	if _defaultLogger == nil {
		return
	}

	_defaultLogger.Write(WarningLevel, tag, format, values...)
}

// Logs an error message to the default logger.
func Error(tag string, err error) {
	if _defaultLogger == nil {
		return
	}

	_defaultLogger.Write(ErrorLevel, tag, "%v", err)
}

// Logs a fatal message to the default logger and panics.
func Fatal(tag string, err error) {
	if _defaultLogger == nil {
		return
	}

	_defaultLogger.Write(FatalLevel, tag, "%v", err)
	panic(err)
}

var (
	_defaultLogger Logger
)

type messageData struct {
	level  Level
	tag    string
	format string
	values []interface{}
}
