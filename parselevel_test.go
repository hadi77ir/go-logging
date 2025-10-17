package logging

import "testing"

func TestParseLevel_Valid(t *testing.T) {
	cases := map[string]Level{
		"panic":   PanicLevel,
		"fatal":   FatalLevel,
		"error":   ErrorLevel,
		"warn":    WarnLevel,
		"warning": WarnLevel,
		"info":    InfoLevel,
		"debug":   DebugLevel,
		"trace":   TraceLevel,
	}
	for in, want := range cases {
		got, err := ParseLevel(in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", in, err)
		}
		if got != want {
			t.Fatalf("%q => %v, want %v", in, got, want)
		}
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	if _, err := ParseLevel("nope"); err == nil {
		t.Fatalf("expected error for invalid level")
	}
}
