package global

import (
	"github.com/op/go-logging"

	"common/cmodel"
)

var Logger cmodel.Logger = nil

type GoLoggingLogger struct {
	Backend logging.LeveledBackend
	Logger  *logging.Logger
}

func (l *GoLoggingLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}
func (l *GoLoggingLogger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *GoLoggingLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}
func (l *GoLoggingLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

func (l *GoLoggingLogger) Warning(args ...interface{}) {
	l.Logger.Warning(args...)
}
func (l *GoLoggingLogger) Warningf(format string, args ...interface{}) {
	l.Logger.Warningf(format, args...)
}

func (l *GoLoggingLogger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}
func (l *GoLoggingLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *GoLoggingLogger) SetLevel(lv int) {
	l.Backend.SetLevel(logging.Level(lv), SERVERNAME)
}
func (l *GoLoggingLogger) GetLevel() int {
	return int(l.Backend.GetLevel(SERVERNAME))
}
