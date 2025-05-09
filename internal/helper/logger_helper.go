package helper

import "github.com/sirupsen/logrus"

type Logger interface {
	Errorf(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Print(args ...interface{})
}

type logger struct {
	log *logrus.Logger
}

func NewLogger(log *logrus.Logger) Logger {
	return &logger{
		log: log,
	}
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
	return
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
	return
}

func (l *logger) Print(args ...interface{}) {
	l.log.Print(args...)
	return
}
