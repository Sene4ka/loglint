package zap

// Logger — заглушка
type Logger struct{}

// NewProduction — заглушка конструктора
func NewProduction() (*Logger, error) { return &Logger{}, nil }

// NewDevelopment — заглушка конструктора
func NewDevelopment() (*Logger, error) { return &Logger{}, nil }

// Методы логгера
func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}

// Field — заглушка для поля логгера
type Field struct{}

// Конструкторы полей (для тестов)
func String(key string, val string) Field   { return Field{} }
func Int(key string, val int) Field         { return Field{} }
func Int64(key string, val int64) Field     { return Field{} }
func Float64(key string, val float64) Field { return Field{} }
func Bool(key string, val bool) Field       { return Field{} }
func Error(err error) Field                 { return Field{} }
func Any(key string, val interface{}) Field { return Field{} }
func Skip() Field                           { return Field{} }

// Zapcore — заглушка для совместимости
type zapcore struct{}

// LevelEnabler — интерфейс для уровня логирования
type LevelEnabler interface {
	Enabled(Level) bool
}

// Level — уровень логирования
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)
