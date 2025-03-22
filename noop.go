package logging

import "os"

type NoOpLogger int

func (n NoOpLogger) Log(level Level, args ...interface{}) {
	// no-op
	if level == FatalLevel {
		os.Exit(1)
	}
}

func (n NoOpLogger) Logger() Logger {
	return n
}

func (n NoOpLogger) WithFields(fields Fields) Logger {
	return n
}

var _ Logger = NoOpLogger(0)
