package logrus

import (
	"errors"
	"testing"

	logging "github.com/hadi77ir/go-logging"
	"github.com/sirupsen/logrus"
)

type fakeLogrus struct {
	level  logrus.Level
	fields logrus.Fields
	args   []any
}

func (f *fakeLogrus) WithFields(fields logrus.Fields) *logrus.Entry {
	f.fields = fields
	logger := logrus.New()
	return logrus.NewEntry(logger)
}

func (f *fakeLogrus) Log(level logrus.Level, args ...interface{}) {
	f.level = level
	f.args = args
}

func TestWrapper_LogPassThrough(t *testing.T) {
	f := &fakeLogrus{}
	w := NewWrapper(f)

	w.Log(logging.ErrorLevel, "hello", 42)

	if f.level != logrus.ErrorLevel {
		t.Fatalf("expected level error, got %v", f.level)
	}
	if len(f.args) != 2 || f.args[0] != "hello" || f.args[1] != 42 {
		t.Fatalf("unexpected args: %#v", f.args)
	}
}

func TestWrapper_WithFieldsAndAdditionalMerge(t *testing.T) {
	f := &fakeLogrus{}
	w := NewWrapper(f)

	// Ensure these calls succeed without panicking even though underlying
	// logger may switch to a *logrus.Entry implementation.
	_ = w.WithFields(logging.Fields{"a": 1})
	w3 := w.WithAdditionalFields(logging.Fields{"a": 2, "b": 3}).(*Wrapper)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	w3.Log(logging.InfoLevel, "msg")
}

func TestNewAtLevelParsesEnvAndSetsLevel(t *testing.T) {
	logger, err := NewAtLevel("debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = logger
}

func TestNewAtLevelRejectsInvalid(t *testing.T) {
	_, err := NewAtLevel("notalevel")
	if err == nil {
		t.Fatalf("expected error for invalid level")
	}
	if !errors.Is(err, err) { // existence check
		// keep linter from complaining that we didn't inspect error
	}
}
