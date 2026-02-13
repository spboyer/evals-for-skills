package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestFailureError(t *testing.T) {
	err := &TestFailureError{
		Message: "benchmark completed with 2 failed and 1 error(s)",
	}

	assert.Equal(t, "benchmark completed with 2 failed and 1 error(s)", err.Error())
}

func TestErrorTypeDetection(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantType string
	}{
		{
			name:     "TestFailureError",
			err:      &TestFailureError{Message: "test failure"},
			wantType: "TestFailureError",
		},
		{
			name:     "regular error",
			err:      errors.New("config error"),
			wantType: "other",
		},
		{
			name:     "wrapped TestFailureError",
			err:      errors.Join(&TestFailureError{Message: "test failure"}, errors.New("additional context")),
			wantType: "TestFailureError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testFailureErr *TestFailureError
			isTestFailure := errors.As(tt.err, &testFailureErr)

			if tt.wantType == "TestFailureError" {
				assert.True(t, isTestFailure, "expected error to be detected as TestFailureError")
			} else {
				assert.False(t, isTestFailure, "expected error NOT to be detected as TestFailureError")
			}
		})
	}
}
