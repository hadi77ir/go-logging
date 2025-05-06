package logrus

import (
	"fmt"
	"github.com/hadi77ir/go-logging"
	"github.com/sirupsen/logrus"
	"os"
)

var _ logging.Logger = &Wrapper{}

type LogrusInterface interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	Log(level logrus.Level, args ...interface{})
}

type Wrapper struct {
	logger LogrusInterface
}

func (c *Wrapper) Logger() logging.Logger {
	if logger, ok := c.logger.(*logrus.Entry); ok {
		return NewWrapper(logger.Logger)
	}
	return c
}

func NewWrapper(logger LogrusInterface) logging.Logger {
	return &Wrapper{logger: logger}
}

func (c *Wrapper) Log(level logging.Level, args ...interface{}) {
	if c.logger != nil {
		c.logger.Log(logrus.Level(level), args)
	}
	// failsafe
	if level == logging.FatalLevel {
		os.Exit(1)
	}
	if level == logging.PanicLevel {
		panic(fmt.Sprint(args...))
	}
}
func (c *Wrapper) WithFields(fields logging.Fields) logging.Logger {
	return NewWrapper(c.logger.WithFields(logrus.Fields(fields)))
}
