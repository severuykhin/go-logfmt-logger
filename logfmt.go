package logfmt

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	L_DEBUG Level = 1
	L_INFO  Level = 2
	L_WARN  Level = 3
	L_ERROR Level = 4
	L_FATAL Level = 5
)

const (
	fieldNameDateTime = "datetime"
	fieldNameMessage  = "message"
	fieldNameLevel    = "level"

	// Optional fieldnames
	fieldNameAppName = "appName"

	fieldValueError = "ERROR"
	fieldValueDebug = "DEBUG"
	fieldValueInfo  = "INFO"
	fieldValueFatal = "FATAL"
	fieldValueWarn  = "WARN"
)

type logger struct {
	output              io.Writer
	verbosityLevel      Level
	levelToTextValueMap map[Level]string
	config              config
}

type config struct {
	AppName   string
	FatalHook func() // Hook on Fatal error level
}

type optFunc = func(c *config)

func WithAppName(appName string) optFunc {
	return func(c *config) {
		c.AppName = appName
	}
}

func WithFatalHook(f func()) optFunc {
	return func(c *config) {
		c.FatalHook = f
	}
}

func New(writer io.Writer, verbosityLevel Level, opts ...optFunc) *logger {

	l := &logger{
		output:         writer,
		verbosityLevel: verbosityLevel,
		levelToTextValueMap: map[Level]string{
			L_DEBUG: fieldValueDebug,
			L_INFO:  fieldValueInfo,
			L_WARN:  fieldValueWarn,
			L_ERROR: fieldValueError,
			L_FATAL: fieldValueFatal,
		},
	}

	var cfg config

	for _, optF := range opts {
		optF(&cfg)
	}

	l.config = cfg

	return l
}

// Useful or important information about the operation of the application
func (l *logger) Info(message any, context ...any) {
	l.log(L_INFO, message, context...)
}

// Additional information about the operation of the application, which may help in identifying errors
func (l *logger) Debug(message any, context ...any) {
	l.log(L_DEBUG, message, context...)
}

// Errors that you can pay attention to, but which do not violate the logic of the application
func (l *logger) Warn(message any, context ...any) {
	l.log(L_WARN, message, context...)
}

// A common error in the process of running an application that needs lighting
func (l *logger) Error(message any, context ...any) {
	l.log(L_ERROR, message, context...)
}

// An error in which further work of applications does not make sense
func (l *logger) Fatal(message any, context ...any) {
	l.log(L_FATAL, message, context...)
}

func (l *logger) log(level Level, message any, context ...any) {

	if level < l.verbosityLevel {
		return
	}

	msg := fmt.Sprintf("\"%v\"", message)

	go func() {
		levelTextValue := l.levelToTextValueMap[level]
		dateTime := time.Now().Format(time.RFC3339)

		stringBuilder := newStringBuilder()

		minLength := len(fieldNameDateTime) +
			len(fieldNameLevel) +
			len(fieldNameMessage) +
			len(dateTime) +
			len(levelTextValue) +
			len([]rune(msg)) +
			3 + 3 // per 1 "=" for each key-val pair + per 1 whitespace for each key-val pair

		stringBuilder.Grow(minLength)

		keyValueSequence := []any{
			fieldNameDateTime, dateTime,
			fieldNameLevel, levelTextValue,
			fieldNameMessage, msg,
			fieldNameAppName, l.config.AppName,
		}

		keyValueSequence = append(keyValueSequence, context...)

		io.WriteString(l.output, stringBuilder.StringFrom(keyValueSequence))

		if level == L_FATAL {
			if l.config.FatalHook != nil {
				l.config.FatalHook()
			} else {
				os.Exit(1)
			}
		}
	}()
}
