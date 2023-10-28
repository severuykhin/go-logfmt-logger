package logfmt

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	l.log(L_INFO, fmt.Sprint(message), context...)
}

// Additional information about the operation of the application, which may help in identifying errors
func (l *logger) Debug(message any, context ...any) {
	l.log(L_DEBUG, fmt.Sprint(message), context...)
}

// Errors that you can pay attention to, but which do not violate the logic of the application
func (l *logger) Warn(message any, context ...any) {
	l.log(L_WARN, fmt.Sprint(message), context...)
}

// A common error in the process of running an application that needs lighting
func (l *logger) Error(message any, context ...any) {
	l.log(L_ERROR, fmt.Sprint(message), context...)
}

// An error in which further work of applications does not make sense
func (l *logger) Fatal(message any, context ...any) {
	l.log(L_FATAL, fmt.Sprint(message), context...)
}

func (l *logger) log(level Level, message string, context ...any) {

	if level < l.verbosityLevel {
		return
	}

	levelTextValue := l.levelToTextValueMap[level]

	// @TODO - refactor string build algorythm
	go func() {
		dateTime := time.Now().Format(time.RFC3339)

		var msgBuilder strings.Builder

		minLength := len(fieldNameMessage) +
			len(fieldNameLevel) +
			len(fieldNameMessage) +
			len(dateTime) +
			len(levelTextValue) +
			len([]rune(message)) +
			3 + 3 // per 1 "=" for each key-val pair + per 1 whitespace for each key-val pair

		msgBuilder.Grow(minLength)

		// msgBuilder.Grow(50)

		l.addKeyValPair(&msgBuilder, fieldNameDateTime, dateTime, false)
		l.addKeyValPair(&msgBuilder, fieldNameLevel, levelTextValue, false)
		l.addKeyValPair(&msgBuilder, fieldNameMessage, message, true)

		if l.config.AppName != "" {
			l.addKeyValPair(&msgBuilder, fieldNameAppName, l.config.AppName, false)
		}

		l.addContextValues(&msgBuilder, context...)

		msgBuilder.WriteString("\n")

		io.WriteString(l.output, msgBuilder.String())

		if level == L_FATAL {
			if l.config.FatalHook != nil {
				l.config.FatalHook()
			} else {
				os.Exit(1)
			}
		}
	}()
}

func (l *logger) addKeyValPair(msgBuilder *strings.Builder, key string, value string, escape bool) {
	msgBuilder.WriteString(key)
	msgBuilder.WriteString("=")
	if escape {
		msgBuilder.WriteString("\"")
	}
	msgBuilder.WriteString(value)
	if escape {
		msgBuilder.WriteString("\"")
	}
	msgBuilder.WriteString(" ")
}

func (l *logger) addContextValues(msgBuilder *strings.Builder, context ...interface{}) {

	contextLength := len(context)

	if contextLength == 0 {
		return
	}

	for i := 0; i < contextLength; i++ {
		if i > 0 && i%2 != 0 {
			continue
		}

		if i == contextLength-1 {
			l.addKeyValPair(msgBuilder, l.valueToString(context[i]), "", false)
		} else {
			l.addKeyValPair(msgBuilder, l.valueToString(context[i]), l.valueToString(context[i+1]), false)
		}

	}
}

func (l *logger) valueToString(param interface{}) string {
	switch v := param.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "unknowntype"
	}
}
