package slog

import "errors"

// A global variable so that log functions can be directly accessed
var log Logger

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	// Debug has verbose message
	Debug = "debug"
	// Info is default log level
	Info = "info"
	// Warn is for logging messages about possible issues
	Warn = "warn"
	// Error is for logging errors
	Error = "error"
	// Fatal is for logging fatal messages. The system shutdown after logging the message.
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Debugw(message string, args ...interface{})

	Infow(message string, args ...interface{})

	Warnw(message string, args ...interface{})

	Errorw(message string, args ...interface{})

	Fatalw(message string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

func GetInstance() Logger {
	return log
}

// NewLogger returns an instance of logger
func NewLogger(logLevel string) {
	// init logger instance
	logConfig := Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      logLevel,
	}
	// set default to zap
	log = newZapLogger(logConfig)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Debugw(message string, args ...interface{}) {
	log.Debugw(message, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Infow(message string, args ...interface{}) {
	log.Infow(message, args...)
}

func Warnw(message string, args ...interface{}) {
	log.Warnw(message, args...)
}

func Errorw(message string, args ...interface{}) {
	log.Errorw(message, args...)
}

// Fatalw can't test this part because after function called os.Exit will be triggered after that
func Fatalw(message string, args ...interface{}) {
	log.Fatalw(message, args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}
