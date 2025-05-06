package logging

import (
	"fmt"
	"strings"
)

// ParseLevel takes a string level and returns the Logrus log level constant. Borrowed from logrus.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

type LimitedLogger struct {
	level  Level
	logger Logger
}

func (l *LimitedLogger) Log(level Level, args ...interface{}) {
	if level < l.level {
		return
	}
	l.logger.Log(l.level, args...)
}

func (l *LimitedLogger) WithFields(fields Fields) Logger {
	return &LimitedLogger{logger: l.logger.WithFields(fields), level: l.level}
}

func (l *LimitedLogger) Logger() Logger {
	return &LimitedLogger{logger: l.logger.Logger(), level: l.level}
}

var _ Logger = &LimitedLogger{}

func Limit(logger Logger, level Level) Logger {
	return &LimitedLogger{logger: logger.Logger(), level: level}
}
