// Package logger represents a generic logging interface
// Taken from: github.com/jfeng45/servicetmpl
package logger

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// var Log Logger

type Logger interface {
	Fatal(args ...interface{})
	Fatalf(tmpl string, args ...interface{})
	Fatalw(msg string, err interface{}, args ...interface{})

	Error(args ...interface{})
	Errorf(tmpl string, args ...interface{})
	Errorw(msg string, err interface{}, args ...interface{})

	Warn(args ...interface{})
	Warnf(tmpl string, args ...interface{})
	Warnw(msg string, args ...interface{})

	Info(args ...interface{})
	Infof(tmpl string, args ...interface{})
	Infow(msg string, args ...interface{})

	Debug(args ...interface{})
	Debugf(tmpl string, args ...interface{})
	Debugw(msg string, args ...interface{})
}

// SetLogger is the setter for log variable, it should be the only way to assign value to log
func SetLogger(newLogger Logger) Logger {
	return newLogger
}

func RegisterLog(level, appMode string) (Logger, error) {
	zLogger, err := initLog(level, appMode)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterLog")
	}
	defer func() {
		tempErr := zLogger.Sync()
		if tempErr != nil {
			log.Println(err)
		}
	}()
	zSugarlog := zLogger.Sugar()

	// This is for loggerWrapper implementation
	return SetLogger(&loggerWrapper{zSugarlog}), nil
	// SetLogger(zSugarlog)
}

func initLog(level string, appMode string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	
	switch appMode {
	case "dev":
		logger, err = zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			return nil, errors.Wrap(err, "zap.NewDevelopment()")
		}
	case "prod":
		logger, err = zap.NewProduction(zap.AddCallerSkip(1))
		if err != nil {
			return nil, errors.Wrap(err, "zap.NewProduction()")
		}
	default:
		return nil, fmt.Errorf("environment variable APP_MODE has unknown value: %s. Exiting", appMode)
	}
	return logger, nil
}

type loggerWrapper struct {
	sl *zap.SugaredLogger
}

// Fatal is for Fatal
func (lg *loggerWrapper) Fatal(args ...interface{}) {
	lg.sl.Fatal(args...)
}

// Fatalf is for Fatalf
func (lg *loggerWrapper) Fatalf(tmpl string, args ...interface{}) {
	lg.sl.Fatalf(tmpl, args...)
}

// Fatalw is for Fatalw
func (lg *loggerWrapper) Fatalw(msg string, err interface{}, args ...interface{}) {
	kvs := make([]interface{}, 0, len(args)+2) // nolint
	kvs = append(kvs, "error", err)
	kvs = append(kvs, args...)
	lg.sl.Fatalw(msg, kvs...)
}

// Error is for Error
func (lg *loggerWrapper) Error(args ...interface{}) {
	lg.sl.Error(args...)
}

// Errorf is for Errorf
func (lg *loggerWrapper) Errorf(tmpl string, args ...interface{}) {
	lg.sl.Errorf(tmpl, args...)
}

// Errorw is for Errorw
func (lg *loggerWrapper) Errorw(msg string, err interface{}, args ...interface{}) {
	kvs := make([]interface{}, 0, len(args)+2) // nolint
	kvs = append(kvs, "error", err)
	kvs = append(kvs, args...)
	lg.sl.Errorw(msg, kvs...)
}

// Warn is for Warn
func (lg *loggerWrapper) Warn(args ...interface{}) {
	lg.sl.Warn(args...)
}

// Warnf is for Warnf
func (lg *loggerWrapper) Warnf(tmpl string, args ...interface{}) {
	lg.sl.Warnf(tmpl, args...)
}

// Warnw is for Warnw
func (lg *loggerWrapper) Warnw(msg string, args ...interface{}) {
	lg.sl.Warnw(msg, args...)
}

// Info is for Info
func (lg *loggerWrapper) Info(args ...interface{}) {
	lg.sl.Info(args...)
}

// Infof is for Infof
func (lg *loggerWrapper) Infof(tmpl string, args ...interface{}) {
	lg.sl.Infof(tmpl, args...)
}

// Infow is for Infow
func (lg *loggerWrapper) Infow(msg string, args ...interface{}) {
	lg.sl.Infow(msg, args...)
}

// Debug is for Debug
func (lg *loggerWrapper) Debug(args ...interface{}) {
	lg.sl.Debug(args...)
}

// Debugf is for Debugf
func (lg *loggerWrapper) Debugf(tmpl string, args ...interface{}) {
	lg.sl.Debugf(tmpl, args...)
}

// Debugw is for Debugw
func (lg *loggerWrapper) Debugw(msg string, args ...interface{}) {
	lg.sl.Debugw(msg, args...)
}

// Sync is for sync
func (lg *loggerWrapper) Sync() {
	err := lg.sl.Sync()
	if err != nil {
		log.Println("Fail to sync zap-logger", err)
	}
}
