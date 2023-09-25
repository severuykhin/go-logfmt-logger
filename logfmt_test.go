package logfmt

import (
	"strings"
	"sync"
	"testing"
	"time"
)

type MockWriter struct {
	buffer []string
	mux    sync.Mutex
}

func (w *MockWriter) Write(p []byte) (int, error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.buffer = append(w.buffer, string(p))
	return len(p), nil
}

func (w *MockWriter) GetBuffer() []string {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.buffer
}

type MockNullWriter struct {
}

func (w *MockNullWriter) Write(p []byte) (int, error) {
	return len(p), nil
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
	logger := New(&w, L_DEBUG, WithFatalHook(func() { /* do nothing */ }))

	logger.Error("some_error_message 1")
	logger.Fatal("some_error_message 2")

	time.Sleep(time.Millisecond * 300)

	for _, line := range w.GetBuffer() {

		t.Run("should set all base fields", func(t *testing.T) {
			if !strings.Contains(line, fieldNameDateTime) {
				t.Fatalf("output string does not contains field: %s", fieldNameDateTime)
			}
			if !strings.Contains(line, fieldNameLevel) {
				t.Fatalf("output string does not contains field: %s", fieldNameLevel)
			}
			if !strings.Contains(line, fieldNameMessage) {
				t.Fatalf("output string does not contains field: %s", fieldNameMessage)
			}
		})

		t.Run("should not set optional fields if they are not specified", func(t *testing.T) {
			if strings.Contains(line, fieldNameAppName) {
				t.Fatalf("output string does not contains field: %s", fieldNameMessage)
			}
		})
	}
}

func TestAppNameOptionalParam(t *testing.T) {

	tearDown := setupTestCase(t)
	defer tearDown(t)

	w := MockWriter{}

	appName := "app_name_test_case"
	logger := New(&w, L_DEBUG, WithAppName(appName), WithFatalHook(func() { /* do nthg */ }))

	logger.Debug("some message")
	logger.Info("some message")
	logger.Warn("some message")
	logger.Error("some message")
	logger.Fatal("some message")

	time.Sleep(time.Millisecond * 300)

	t.Run("should append optional appName param in all cases", func(t *testing.T) {
		for _, line := range w.GetBuffer() {
			if !strings.Contains(line, fieldNameAppName) {
				t.Fatalf("output string does not contains appname field: %s", fieldNameAppName)
			}
			if !strings.Contains(line, appName) {
				t.Fatalf("output string does not contains appname value: %s", appName)
			}
		}
	})
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

	logger.Debug("message", "param1", "value1", "param2", 42, 123, "value3")

	time.Sleep(time.Millisecond * 300)

	if len(w.GetBuffer()) == 0 {
		t.Fatal("buffer does not contains output")
	}

	line := w.GetBuffer()[0]

	if !strings.Contains(line, "param1=value1") {
		t.Fatalf("output string does not contains field: %s, output: %s", "param1", line)
	}

	if !strings.Contains(line, "param2=42") {
		t.Fatalf("output string does not contains field: %s, output: %s", "param2", line)
	}

	if !strings.Contains(line, "123=value3") {
		t.Fatalf("output string does not contains field: %s, output: %s", "value3", line)
	}
}

func BenchmarkWithoutContext(b *testing.B) {
	w := MockNullWriter{}
	logger := New(&w, L_DEBUG, WithFatalHook(func() { /* do nothing */ }))
	for i := 0; i < b.N; i++ {
		logger.Error("test message")
	}
}

func BenchmarkWithContext(b *testing.B) {
	w := MockNullWriter{}
	logger := New(&w, L_DEBUG, WithFatalHook(func() { /* do nothing */ }))
	for i := 0; i < b.N; i++ {
		logger.Error("test message", "code", 123, "key", "value")
	}
}
