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
	L_ERROR Level = 3
	L_FATAL Level = 4
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
			L_DEBUG: "DEBUG",
			L_INFO:  "INFO",
			L_ERROR: "ERROR",
			L_FATAL: "FATAL",
		},
	}
}

// A common error in the process of running an application that needs lighting
func (l *logger) Error(code int, message string) {
	l.log(L_ERROR, code, message)
}

// Useful or important information about the operation of the application
func (l *logger) Info(code int, message string) {
	l.log(L_INFO, code, message)
}

// Additional information about the operation of the application, which may help in identifying errors
func (l *logger) Debug(code int, message string) {
	l.log(L_DEBUG, code, message)
}

// An error in which further work of applications does not make sense
func (l *logger) Fatal(code int, message string) {
	l.log(L_FATAL, code, message)
}

func (l *logger) log(level Level, code int, message string) {

	if level < l.verbosityLevel {
		return
	}

	levelTextValue := l.levelToTextValueMap[level]

	go func() {
		dateTime := time.Now().Format(time.RFC3339)

		var msgTemplate strings.Builder

		msgTemplate.WriteString("datetime=%s ")
		msgTemplate.WriteString("level=%s ")
		msgTemplate.WriteString("code=%s ")
		msgTemplate.WriteString("message=\"%s\" ")
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
