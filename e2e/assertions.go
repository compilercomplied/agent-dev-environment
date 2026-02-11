package e2e

import (
	"errors"
	"testing"
)

// AssertError provides a one-line assertion for expected API errors.
func AssertError(t *testing.T, err error, expectedStatus int, expectedMessage string) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected APIError, got %T: %v", err, err)
	}

	if apiErr.Status != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, apiErr.Status)
	}

	if expectedMessage != "" && apiErr.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, apiErr.Message)
	}
}
