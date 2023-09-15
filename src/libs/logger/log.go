package logger

import (
	"fmt"
	"runtime/debug"
)

func (l *Logger) Log(message interface{}, options ...*Options) {
	msg := []interface{}{message}

	if len(options) > 0 && options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	l.stdout.Println(msg...)
}
