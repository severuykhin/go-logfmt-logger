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
	fieldNameCode     = "code"
	fieldNameMessage  = "message"
	fieldNameLevel    = "level"

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
}

func New(writer io.Writer, verbosityLevel Level) *logger {
	return &logger{
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
}

// Useful or important information about the operation of the application
func (l *logger) Info(code int, message string, context ...interface{}) {
	l.log(L_INFO, code, message, context...)
}

// Additional information about the operation of the application, which may help in identifying errors
func (l *logger) Debug(code int, message string, context ...interface{}) {
	l.log(L_DEBUG, code, message, context...)
}

// Errors that you can pay attention to, but which do not violate the logic of the application
func (l *logger) Warn(code int, message string, context ...interface{}) {
	l.log(L_WARN, code, message, context...)
}

// A common error in the process of running an application that needs lighting
func (l *logger) Error(code int, message string, context ...interface{}) {
	l.log(L_ERROR, code, message, context...)
}

// An error in which further work of applications does not make sense
func (l *logger) Fatal(code int, message string, context ...interface{}) {
	l.log(L_FATAL, code, message, context...)
}

func (l *logger) log(level Level, code int, message string, context ...interface{}) {

	if level < l.verbosityLevel {
		return
	}

	levelTextValue := l.levelToTextValueMap[level]

	go func() {
		dateTime := time.Now().Format(time.RFC3339)

		var msgTemplate strings.Builder

		msgTemplate.WriteString(fieldNameDateTime + "=%s ")
		msgTemplate.WriteString(fieldNameLevel + "=%s ")
		msgTemplate.WriteString(fieldNameCode + "=%s ")
		msgTemplate.WriteString(fieldNameMessage + "=\"%s\" ")

		l.addContextValues(&msgTemplate, context...)

		msgTemplate.WriteString("\n")

		msg := fmt.Sprintf(
			msgTemplate.String(),
			dateTime,
			levelTextValue,
			strconv.Itoa(code),
			message)

		io.WriteString(l.output, msg)

		if level == L_FATAL {
			os.Exit(1)
		}
	}()
}

func (l *logger) addContextValues(strBuilder *strings.Builder, context ...interface{}) {

	contextLength := len(context)

	if contextLength == 0 {
		return
	}

	for i := 0; i < contextLength; i++ {
		if i > 0 && i%2 != 0 {
			continue
		}

		if i == contextLength-1 {
			strBuilder.WriteString(
				l.getFormattedParam(context[i]) + "=" + " ",
			)
		} else {
			strBuilder.WriteString(
				l.getFormattedParam(context[i]) + "=" + l.getFormattedValue(context[i+1]) + " ",
			)
		}

	}
}

func (l *logger) getFormattedParam(param interface{}) string {
	switch v := param.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return "unknowntype"
	}
}

func (l *logger) getFormattedValue(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return fmt.Sprintf("\"%s\"", v)
	default:
		//TODO - add rest types
		return "unknowntype"
	}
}
