package run

import (
	"agent-dev-environment/e2e"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	ls_models "agent-dev-environment/src/api/v1/filesystem/ls"
	run_models "agent-dev-environment/src/api/v1/shell/run"
	"strings"
	"testing"
)

func TestRunShell_GitClone(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	repoDir := "/tmp/agent-dev-environment"

	// Cleanup before test in case it exists from a previous run
	client.DeleteFile(delete_models.Request{
		Path:      repoDir,
		Recursive: true,
	})

	req := run_models.Request{
		Command: "sh",
		Args:    []string{"-c", "cd /tmp && git clone https://github.com/compilercomplied/agent-dev-environment.git"},
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.RunShell(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error from RunShell, got %v", err)
	}

	// Verify the directory exists by listing it
	lsResp, err := client.ListFiles(ls_models.Request{
		Path: "/tmp",
	})
	if err != nil {
		t.Fatalf("Expected no error from ListFiles, got %v", err)
	}

	if !strings.Contains(lsResp.CommandOutput, "agent-dev-environment") {
		t.Errorf("Expected 'agent-dev-environment' in /tmp listing, but got: %s. RunShell output: %s", lsResp.CommandOutput, resp.CommandOutput)
	}

	// Cleanup after test
	client.DeleteFile(delete_models.Request{
		Path:      repoDir,
		Recursive: true,
	})
}
