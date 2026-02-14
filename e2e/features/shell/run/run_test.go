package run

import (
	. "agent-dev-environment/e2e"
	run_models "agent-dev-environment/src/api/v1/shell/run"
	"strings"
	"testing"
)

func TestRunShell_AllowedCommand_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := run_models.Request{
		Command: "ls",
		Args:    []string{"."},
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.RunShell(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.CommandOutput == "" {
		t.Error("Expected command output, got empty string")
	}
}

func TestRunShell_CurlLocalhost_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := run_models.Request{
		Command: "curl",
		Args:    []string{"-s", "http://localhost:8080/health"},
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.RunShell(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(resp.CommandOutput, "OK") {
		t.Errorf("Expected output to contain 'OK', got %q", resp.CommandOutput)
	}
}

func TestRunShell_CurlExternal_Restricted(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := run_models.Request{
		Command: "curl",
		Args:    []string{"http://google.com"},
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.RunShell(req)

	// ------------------------------------ Assert -------------------------------------
	if err == nil {
		t.Fatal("Expected error for external curl, got none")
	}
	if !strings.Contains(err.Error(), "curl is restricted to localhost targets due to security reasons") {
		t.Errorf("Expected restriction error message, got: %v", err)
	}
}

func TestRunShell_UnauthorizedCommand_Restricted(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := run_models.Request{
		Command: "apt-get",
		Args:    []string{"install", "vim"},
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.RunShell(req)

	// ------------------------------------ Assert -------------------------------------
	if err == nil {
		t.Fatal("Expected error for unauthorized command, got none")
	}
	if !strings.Contains(err.Error(), "command 'apt-get' is not allowed") {
		t.Errorf("Expected 'not allowed' error message, got: %v", err)
	}
}
