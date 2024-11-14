package simplelog

import "testing"

func TestSetLogLevel(t *testing.T) {
	l := NewLogger(INFO)
	l.SetLevel(DEBUG)
	if l.level != DEBUG {
		t.Errorf("expected log level %d, got %d", DEBUG, l.level)
	}

	t.Logf("log level set to %d", l.level)
}

func TestMessage(t *testing.T) {
	l := NewLogger(DEBUG)
	l.Debug("debug message")
	Info("info message")
}
