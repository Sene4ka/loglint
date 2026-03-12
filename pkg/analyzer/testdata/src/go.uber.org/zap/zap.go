package zap

type Logger struct{}

func (*Logger) Info(msg string, fields ...interface{})  {}
func (*Logger) Error(msg string, fields ...interface{}) {}
func (*Logger) Debug(msg string, fields ...interface{}) {}
func (*Logger) Warn(msg string, fields ...interface{})  {}
