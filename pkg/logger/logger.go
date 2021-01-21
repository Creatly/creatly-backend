package logger

type Logger interface {
	Debug(msg string, params map[string]interface{})
	Info(msg string, params map[string]interface{})
	Warn(msg string, params map[string]interface{})
	Error(msg string, params map[string]interface{})
}
