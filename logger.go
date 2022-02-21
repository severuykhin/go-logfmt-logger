package logfmt

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	L_DEBUG = "debug"
	L_INFO  = "info"
	L_ERROR = "error"
	L_FATAL = "fatal"
)

type logger struct {
	output io.Writer
}

func New(writer io.Writer) *logger {
	return &logger{
		output: writer,
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

func (l *logger) log(level string, code int, message string) {
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
			level,
			strconv.Itoa(code),
			message)

		io.WriteString(l.output, msg)

		if level == L_FATAL {
			os.Exit(1)
		}
	}()
}
