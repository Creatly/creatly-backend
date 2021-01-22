package logger

import (
	"github.com/sirupsen/logrus"
)

// TODO create general interface with generic fields

//func Init() error {
//	var err error
//	logger, err := zap.NewProduction()
//	if err != nil {
//		return err
//	}
//
//	zap.ReplaceGlobals(logger)
//
//	return err
//}

func Debug(msg... interface{}) {
	//zap.S().Debug(msg)
	logrus.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	//zap.S().Debugf(format, args)
	logrus.Debugf(format, args)
}

func Info(msg... interface{}) {
	//zap.S().Info(msg)
	logrus.Info(msg)
}

func Infof(format string, args ...interface{}) {
	//zap.S().Infof(format, args)
	logrus.Infof(format, args)
}

func Warn(msg... interface{}) {
	//zap.S().Warn(msg)
	logrus.Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	//zap.S().Warnf(format, args)
	logrus.Warnf(format, args)
}

func Error(msg... interface{}) {
	logrus.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	//zap.S().Errorf(format, args)
	logrus.Errorf(format, args)
}
