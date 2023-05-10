package logger

import (
	"go.uber.org/zap"
)

// Logger is a utility struct for logging data in an extremely high performance system.
// We can use both Logger and SugarLog for logging. For more information,
// just visit https://godoc.org/go.uber.org/zap
type Logger struct {
	sl *zap.SugaredLogger
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sl.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sl.Fatalf(format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.sl.Fatal(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sl.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sl.Panicf(format, args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.sl.Panic(args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.sl.Info(args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.sl.Infof(format, args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.sl.Info(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.sl.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sl.Debugf(format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.sl.Debug(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sl.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sl.Errorf(format, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.sl.Error(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sl.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sl.Infof(format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.sl.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sl.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sl.Warnf(format, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.sl.Warn(args...)
}
