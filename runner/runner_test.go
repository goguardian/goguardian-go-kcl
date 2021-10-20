package runner

import (
	"io"
	"log"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

type mockIOWriter struct {
	logChan chan []byte
}

func (m *mockIOWriter) Write(p []byte) (int, error) {
	m.logChan <- p
	return len(p), nil
}

func TestPipeToLogger_SendsOutputToLogger(t *testing.T) {
	// setup
	someLogOutput := "foo bar\n"
	logChan := make(chan []byte)
	mockWriter := &mockIOWriter{logChan: logChan}
	logger := log.New(mockWriter, "", 0)

	getPipe := func() (io.ReadCloser, error) {
		r := io.NopCloser(strings.NewReader(someLogOutput))
		return r, nil
	}

	// test
	err := pipeToLogger(logger, getPipe)
	if err != nil {
		t.Error(err)
	}

	// validate
	actualLogLine := <-logChan
	if string(actualLogLine) != someLogOutput {
		t.Errorf("expected '%s' log but got '%s'", someLogOutput, actualLogLine)
	}
}

func TestPipeToLogger_ReturnsErrorWhenGetPipeReturnsError(t *testing.T) {
	// setup
	someError := errors.New("some error")
	logger := log.New(nil, "", 0)

	getPipe := func() (io.ReadCloser, error) {
		return nil, someError
	}

	// test
	err := pipeToLogger(logger, getPipe)

	// validate
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}
