package logging

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

func (l *LimitedLogger) WithAdditionalFields(fields Fields) Logger {
	return &LimitedLogger{logger: l.logger.WithAdditionalFields(fields), level: l.level}
}

func (l *LimitedLogger) Logger() Logger {
	return &LimitedLogger{logger: l.logger.Logger(), level: l.level}
}

var _ Logger = &LimitedLogger{}

func Limit(logger Logger, level Level) Logger {
	return &LimitedLogger{logger: logger.Logger(), level: level}
}
