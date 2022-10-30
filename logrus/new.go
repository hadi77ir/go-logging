package logrus

import (
	"github.com/hadi77ir/go-logging"
	"github.com/sirupsen/logrus"
	"os"
)

// New creates a new logger instance.
// it should be used by main() function and also each command's Run() function
// to initialize the logging functionality.
func New(tag string) (logging.Logger, error) {
	return NewAtLevel(tag, os.Getenv("LOG_LEVEL"))
}

func NewAtLevel(tag string, levelStr string) (logging.Logger, error) {
	logger := logrus.New()
	logLevel := logrus.InfoLevel
	if levelStr != "" {
		var err error
		logLevel, err = logrus.ParseLevel(levelStr)
		if err != nil {
			return nil, err
		}
	}

	logger.SetLevel(logLevel)

	return NewLogrusWrapper(logger.WithField("command", tag)), nil
}
