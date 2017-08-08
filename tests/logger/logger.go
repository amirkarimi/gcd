package logger

type Logger struct{}

func (l Logger) Info(message ...interface{}) {}

func (l Logger) Error(message ...interface{}) {}

func (l Logger) OK(message ...interface{}) {}

func (l Logger) Exit(code int, message ...interface{}) {}

func (l Logger) Skip(message ...interface{}) {}
