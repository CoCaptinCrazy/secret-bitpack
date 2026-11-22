package logging

import (
	"fmt"
	"io"
)

type Logger struct {
	output io.Writer
}

func New(output io.Writer) *Logger {
	return &Logger{output: output}
}

func (l *Logger) Info(format string, args ...any) {
	fmt.Fprintf(l.output, "[info] "+format+"\n", args...)
}

func (l *Logger) Error(format string, args ...any) {
	fmt.Fprintf(l.output, "[error] "+format+"\n", args...)
}
