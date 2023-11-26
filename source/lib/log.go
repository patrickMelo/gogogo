package lib

import (
	"fmt"
	"time"
)

// Enables verbose logging.
func EnableVerboseLog() {
	logVerboseEnabled = true
}

// Logs a information message to the standard output.
func LogInformation(tag string, format string, values ...interface{}) {
	writeLog('I', tag, format, values...)
}

// Logs verbose messages to the standard output (if verbose log is enabled).
func LogVerbose(tag string, format string, values ...interface{}) {
	if !logVerboseEnabled {
		return
	}

	writeLog('V', tag, format, values...)
}

// Logs a warning message to the standard output.
func LogWarning(tag string, format string, values ...interface{}) {
	writeLog('W', tag, format, values...)
}

// Logs an error message to the standard output.
func LogError(tag string, err error) {
	writeLog('E', tag, "%v", err)
}

// Logs a fatal message to the standard output and stops the service.
func LogFatal(tag string, err error) {
	writeLog('F', tag, "%v", err)
	serviceRunning = false
}

var (
	logVerboseEnabled = false
)

func writeLog(typeId rune, tag string, format string, values ...interface{}) {
	fmt.Printf("[%s] (%c) [%s] %s\n", time.Now().Format("2006-01-02 15:04:05.000000"), typeId, tag, fmt.Sprintf(format, values...))
}
