package logger

import (
	"go.uber.org/zap"
)

// TODO create general interface with generic fields

func Init() error {
	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return err
}

func Debug(msg... interface{}) {
	zap.S().Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	zap.S().Debugf(format, args)
}

func Info(msg... interface{}) {
	zap.S().Info(msg)
}

func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args)
}

func Warn(msg... interface{}) {
	zap.S().Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	zap.S().Warnf(format, args)
}

func Error(msg... interface{}) {
	zap.S().Error(msg)
}

func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args)
}
