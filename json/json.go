package json

import (
	"fmt"
	"io"
	"maps"
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/hadi77ir/go-logging"
)

func levelString(level logging.Level) string {
	switch level {
	case logging.TraceLevel:
		return "trace"
	case logging.DebugLevel:
		return "debug"
	case logging.InfoLevel:
		return "info"
	case logging.WarnLevel:
		return "warning"
	case logging.ErrorLevel:
		return "error"
	case logging.FatalLevel:
		return "fatal"
	case logging.PanicLevel:
		return "panic"
	}
	return "unknown"
}

type jsonWriter struct {
	encoder *json.Encoder
	mutex   sync.Mutex
}

func (w *jsonWriter) Write(level logging.Level, args []any, fields logging.Fields, timestamp time.Time) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	content := logging.Fields{
		"timestamp": timestamp.UTC().Format(time.RFC3339),
		"level":     levelString(level),
		"message":   fmt.Sprint(args...),
	}
	maps.Copy(content, fields)
	_ = w.encoder.Encode(content)
}

func newWriter(writer io.Writer) *jsonWriter {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "")
	return &jsonWriter{
		encoder: encoder,
	}
}

type Logger struct {
	writer *jsonWriter
	fields logging.Fields
}

func (l *Logger) Log(level logging.Level, args ...interface{}) {
	l.writer.Write(level, args, l.fields, time.Now())

	if level == logging.FatalLevel {
		os.Exit(1)
	}
	if level == logging.PanicLevel {
		panic(fmt.Sprint(args...))
	}
}

func (l *Logger) WithFields(fields logging.Fields) logging.Logger {
	return &Logger{
		writer: l.writer,
		fields: fields,
	}
}

func (l *Logger) WithAdditionalFields(fields logging.Fields) logging.Logger {
	merged := fields
	for k, v := range l.fields {
		if _, ok := merged[k]; !ok {
			merged[k] = v
		}
	}
	return l.WithFields(merged)
}

func (l *Logger) Logger() logging.Logger {
	return &Logger{writer: l.writer}
}

func New(writer io.Writer) *Logger {
	return &Logger{writer: newWriter(writer)}
}
func NewStderr() *Logger {
	return New(os.Stderr)
}
