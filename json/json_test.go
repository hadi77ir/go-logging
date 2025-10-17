package json

import (
	"bytes"
	stdjson "encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	logging "github.com/hadi77ir/go-logging"
)

func readLine(buf *bytes.Buffer) (map[string]any, error) {
	line, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := stdjson.Unmarshal(line, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func TestLogger_WritesJSONWithFieldsAndLevels(t *testing.T) {
	var out bytes.Buffer
	l := New(&out)

	l.Log(logging.InfoLevel, "hello", " ", "world")

	obj, err := readLine(&out)
	if err != nil {
		t.Fatalf("failed reading line: %v", err)
	}

	if obj["level"] != "info" {
		t.Fatalf("expected level=info, got %v", obj["level"])
	}

	msg, _ := obj["message"].(string)
	if msg != "hello world" {
		t.Fatalf("expected message 'hello world', got %q", msg)
	}

	// timestamp should be RFC3339 and in UTC
	ts, _ := obj["timestamp"].(string)
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		t.Fatalf("timestamp not RFC3339: %v (value: %q)", err, ts)
	}
}

func TestLogger_WithFieldsAndOverrides(t *testing.T) {
	var out bytes.Buffer
	l := New(&out)

	l2 := l.WithFields(logging.Fields{"a": 1, "message": "override"})
	l2.Log(logging.DebugLevel, "ignored")

	obj, err := readLine(&out)
	if err != nil {
		t.Fatalf("failed reading line: %v", err)
	}

	if obj["a"] != float64(1) { // json numbers become float64
		t.Fatalf("expected field a=1, got %v", obj["a"])
	}

	// user-provided fields override default keys
	if obj["message"] != "override" {
		t.Fatalf("expected message to be overridden, got %v", obj["message"])
	}

	// original logger should remain without fields
	l.Log(logging.InfoLevel, "base")
	obj2, err := readLine(&out)
	if err != nil {
		t.Fatalf("failed reading second line: %v", err)
	}
	if _, ok := obj2["a"]; ok {
		t.Fatalf("unexpected field 'a' on base logger output: %v", obj2)
	}
}

func TestLogger_WithAdditionalFields_Merges(t *testing.T) {
	var out bytes.Buffer
	l := New(&out).WithFields(logging.Fields{"a": 1}).(*Logger)

	l3 := l.WithAdditionalFields(logging.Fields{"b": 2, "a": 99})
	l3.Log(logging.WarnLevel, "payload")

	obj, err := readLine(&out)
	if err != nil {
		t.Fatalf("failed reading line: %v", err)
	}

	if obj["a"] != float64(99) {
		t.Fatalf("expected merged 'a'=99, got %v", obj["a"])
	}
	if obj["b"] != float64(2) {
		t.Fatalf("expected merged 'b'=2, got %v", obj["b"])
	}
}

func TestLogger_PanicLevelPanics(t *testing.T) {
	var out bytes.Buffer
	l := New(&out)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		} else {
			if !strings.Contains(r.(string), "boom") {
				t.Fatalf("unexpected panic message: %v", r)
			}
		}
	}()

	l.Log(logging.PanicLevel, "boom")
}

func TestLogger_FatalLevelExits(t *testing.T) {
	if os.Getenv("JSON_TEST_TRIGGER_FATAL") == "1" {
		var out bytes.Buffer
		l := New(&out)
		l.Log(logging.FatalLevel, "fatal")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run", "TestLogger_FatalLevelExits")
	cmd.Env = append(os.Environ(), "JSON_TEST_TRIGGER_FATAL=1")
	if err := cmd.Run(); err == nil {
		t.Fatalf("expected process to exit with error, got nil")
	}
}
