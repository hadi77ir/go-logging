package logrus

import (
	"fmt"
	"github.com/hadi77ir/go-logging"
	"github.com/sirupsen/logrus"
	"maps"
	"os"
)

var _ logging.Logger = &Wrapper{}

type LogrusInterface interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	Log(level logrus.Level, args ...interface{})
}

type Wrapper struct {
	logger LogrusInterface
	fields logging.Fields
}

func (c *Wrapper) Logger() logging.Logger {
	if logger, ok := c.logger.(*logrus.Entry); ok {
		return NewWrapper(logger.Logger)
	}
	return c
}

func NewWrapper(logger LogrusInterface) logging.Logger {
	return newWrapperWithFields(logger, logging.Fields{})
}

func newWrapperWithFields(logger LogrusInterface, fields logging.Fields) logging.Logger {
	return &Wrapper{logger: logger, fields: fields}
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
	return newWrapperWithFields(c.logger.WithFields(logrus.Fields(fields)), fields)
}
func (c *Wrapper) WithAdditionalFields(fields logging.Fields) logging.Logger {
	// no need to clone, as fields map shouldn't be modified by caller.
	merged := fields
	maps.Copy(merged, c.fields)
	return newWrapperWithFields(c.logger.WithFields(logrus.Fields(merged)), merged)
}
