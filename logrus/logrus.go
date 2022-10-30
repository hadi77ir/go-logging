package logrus

import (
	"github.com/hadi77ir/go-logging"
	"github.com/sirupsen/logrus"
)

var _ logging.Logger = &LogrusWrapper{}

type LogrusInterface interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	Log(level logrus.Level, args ...interface{})
}

type LogrusWrapper struct {
	logger LogrusInterface
}

func (c *LogrusWrapper) Logger() logging.Logger {
	if logger, ok := c.logger.(*logrus.Entry); ok {
		return NewLogrusWrapper(logger.Logger)
	}
	return c
}

func NewLogrusWrapper(logger LogrusInterface) logging.Logger {
	return &LogrusWrapper{logger: logger}
}

func (c *LogrusWrapper) Log(level logging.Level, args ...interface{}) {
	if c.logger != nil {
		c.logger.Log(logrus.Level(level), args)
	}
}
func (c *LogrusWrapper) WithFields(fields logging.Fields) logging.Logger {
	return NewLogrusWrapper(c.logger.WithFields(logrus.Fields(fields)))
}
