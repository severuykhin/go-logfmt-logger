package logfmt

import (
	"strings"
	"testing"
	"time"
)

type MockWriter struct {
	buffer []string
}

func (w *MockWriter) Write(p []byte) (int, error) {
	w.buffer = append(w.buffer, string(p))
	return len(p), nil
}

func (w *MockWriter) GetBuffer() []string {
	return w.buffer
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		// to do
	}
}

func TestOutputIsInCorrectFormat(t *testing.T) {

	tearDown := setupTestCase(t)
	defer tearDown(t)

	w := MockWriter{}
	logger := New(&w, L_DEBUG)

	logger.Error(123, "some_error_message 1")
	logger.Error(456, "some_error_message 2")

	time.Sleep(time.Millisecond * 300)

	for _, line := range w.GetBuffer() {
		if !strings.Contains(line, fieldNameDateTime) {
			t.Fatalf("output string does not contains field: %s", fieldNameDateTime)
		}
		if !strings.Contains(line, fieldNameCode) {
			t.Fatalf("output string does not contains field: %s", fieldNameCode)
		}
		if !strings.Contains(line, fieldNameLevel) {
			t.Fatalf("output string does not contains field: %s", fieldNameLevel)
		}
		if !strings.Contains(line, fieldNameMessage) {
			t.Fatalf("output string does not contains field: %s", fieldNameMessage)
		}
	}
}

// func TestVerbosityLevelDebug(t *testing.T) {

// 	tearDown := setupTestCase(t)
// 	defer tearDown(t)

// 	logger := New(MockWriter{}, L_DEBUG)

// 	logger.Debug(1, "some message")
// 	logger.Info(1, "some message")
// 	logger.Warn(1, "some message")
// 	logger.Error(1, "some message")

// 	time.Sleep(time.Millisecond * 1000)

// }

func TestMessageContextParams(t *testing.T) {
	tearDown := setupTestCase(t)
	defer tearDown(t)

	w := MockWriter{}
	logger := New(&w, L_DEBUG)

	logger.Debug(502, "message", "param1", "value1", "param2", 42, 123, "value3")

	time.Sleep(time.Millisecond * 300)

	if len(w.GetBuffer()) == 0 {
		t.Fatal("buffer does not contains output")
	}

	line := w.GetBuffer()[0]

	if !strings.Contains(line, "param1=\"value1\"") {
		t.Fatalf("output string does not contains field: %s, output: %s", "param1", line)
	}

	if !strings.Contains(line, "param2=42") {
		t.Fatalf("output string does not contains field: %s, output: %s", "param2", line)
	}

	if !strings.Contains(line, "123=\"value3\"") {
		t.Fatalf("output string does not contains field: %s, output: %s", "value3", line)
	}
}
