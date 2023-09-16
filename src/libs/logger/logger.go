package logger

import (
	"log"
	"os"
)

type Logger struct {
	stdout *log.Logger
	stderr *log.Logger
}

type Options struct {
	IsPrintStack bool
	IsExit       bool
	ExitCode     int
}

func New(prefix string) *Logger {
	return &Logger{
		stdout: log.New(os.Stdout, "[LOG]["+prefix+"]", log.Ldate|log.Ltime),
		stderr: log.New(os.Stderr, "[ERROR]["+prefix+"]", log.Ldate|log.Ltime),
	}
}
