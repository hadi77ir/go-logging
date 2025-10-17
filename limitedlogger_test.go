package logging

import "testing"

type captureLogger struct {
	entries []struct {
		level Level
		args  []any
	}
}

func (c *captureLogger) Log(level Level, args ...interface{}) {
	c.entries = append(c.entries, struct {
		level Level
		args  []any
	}{level, args})
}
func (c *captureLogger) WithFields(fields Fields) Logger           { return c }
func (c *captureLogger) WithAdditionalFields(fields Fields) Logger { return c }
func (c *captureLogger) Logger() Logger                            { return c }

func TestLimitedLogger_FiltersByLevel(t *testing.T) {
	base := &captureLogger{}
	l := Limit(base, WarnLevel)

	l.Log(DebugLevel, "debug")
	l.Log(InfoLevel, "info")
	l.Log(WarnLevel, "warn")
	l.Log(ErrorLevel, "error")

	if len(base.entries) != 2 {
		t.Fatalf("expected 2 entries >= warn, got %d", len(base.entries))
	}
	if base.entries[0].level != WarnLevel || base.entries[1].level != ErrorLevel {
		t.Fatalf("unexpected levels: %#v", base.entries)
	}
}

func TestLimitedLogger_PreservesChaining(t *testing.T) {
	base := &captureLogger{}
	l := Limit(base, InfoLevel)

	l2 := l.WithFields(Fields{"a": 1})
	_ = l2.WithAdditionalFields(Fields{"b": 2})

	l2.Log(InfoLevel, "x")
	if len(base.entries) != 1 || base.entries[0].level != InfoLevel {
		t.Fatalf("unexpected entries: %#v", base.entries)
	}
}
