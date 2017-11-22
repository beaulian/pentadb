package log

import (
	"io"
	"os"
	"log"
	"fmt"
)

// log level
const (
	Info = 32
	Warning = 33
	Error = 31
)

type Log struct {
	// standard library
	logger *log.Logger

	colorTemplate string

	infoColor int
	warningColor int
	errorColor int
}


func NewLog(out io.Writer) *Log {
	if out == nil {
		out = os.Stdout
	}
	logger := log.New(out, "", log.Ldate | log.Ltime)
	return &Log {
		logger:         logger,
		colorTemplate:  "0x1b[1;%dm%s%c[0m\n",
		infoColor:      Info,
		warningColor:   Warning,
		errorColor:     Error,
	}
}

func (l *Log) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

// set color template to replace the old
func (l *Log) SetColorTemplate(ct string) {
	l.colorTemplate = ct
}

// set color for each level
func (l *Log) SetInfoColor(color int) {
	l.infoColor = color
}

func (l *Log) SetWarningColor(color int) {
	l.warningColor = color
}

func (l *Log) SetErrorColor(color int) {
	l.errorColor = color
}

func (l *Log) Info(format string, v ...interface{}) {
	text := fmt.Sprintf(format, v)
	l.logger.Printf(l.colorTemplate, 0x1B, l.infoColor, text, 0x1B)
}

func (l *Log) Warning(format string, v ...interface{}) {
	text := fmt.Sprintf(format, v)
	l.logger.Panicf(l.colorTemplate, 0x1B, l.warningColor, text, 0x1B)
}

func (l *Log) Error(format string, v ...interface{}) {
	text := fmt.Sprintf(format, v)
	l.logger.Fatalf(l.colorTemplate, 0x1B, l.errorColor, text, 0x1B)
}