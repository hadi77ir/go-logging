package logging

import (
	"os"
	"os/exec"
	"testing"
)

func TestNoOpLogger_DoesNothing(t *testing.T) {
	var n NoOpLogger
	// should not panic
	n.Log(InfoLevel, "hello")
	if n.Logger() != n {
		t.Fatalf("Logger() should return itself for NoOpLogger")
	}
	if n.WithFields(Fields{"a": 1}) != n {
		t.Fatalf("WithFields should be no-op for NoOpLogger")
	}
	if n.WithAdditionalFields(Fields{"a": 1}) != n {
		t.Fatalf("WithAdditionalFields should be no-op for NoOpLogger")
	}
}

func TestNoOpLogger_FatalExits(t *testing.T) {
	if os.Getenv("NOOP_TEST_TRIGGER_FATAL") == "1" {
		var n NoOpLogger
		n.Log(FatalLevel, "bye")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run", "TestNoOpLogger_FatalExits")
	cmd.Env = append(os.Environ(), "NOOP_TEST_TRIGGER_FATAL=1")
	if err := cmd.Run(); err == nil {
		t.Fatalf("expected process to exit with error, got nil")
	}
}

func TestNoOpLogger_PanicPanics(t *testing.T) {
	var n NoOpLogger
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	n.Log(PanicLevel, "boom")
}
