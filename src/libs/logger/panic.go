package logger

func (l *Logger) Panic(message interface{}, options ...Options) {
	l.stderr.Panicln(message)
}
