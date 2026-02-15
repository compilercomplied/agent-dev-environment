package reload_env

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	reload_models "agent-dev-environment/src/api/v1/shell/reload_env"
	run_models "agent-dev-environment/src/api/v1/shell/run"
	"strings"
	"testing"
)

func TestReloadEnv_MockedDotEnv(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	testVarName := "E2E_MOCKED_VAR"
	testVarValue := "mock-confirmed"
	envPath := "test-reload.env"
	
	_, err := client.CreateFile(create_models.Request{
		Path:    envPath,
		Content: testVarName + "='" + testVarValue + "'\n",
	})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Failed to arrange: could not create env file: %v", err)
	}

	defer func() {
		client.DeleteFile(delete_models.Request{Path: envPath})
	}()

	// -------------------------------------- Act --------------------------------------
	_, err = client.ReloadEnv(reload_models.Request{Path: envPath})
	if err != nil {
		t.Fatalf("Act failed: ReloadEnv returned error: %v", err)
	}

	// ------------------------------------ Assert -------------------------------------
	resp, err := client.RunShell(run_models.Request{
		Command: "env",
		Args:    []string{},
	})

	if err != nil {
		t.Fatalf("Failed to assert: shell command failed: %v", err)
	}

	expectedLine := testVarName + "=" + testVarValue
	if !strings.Contains(resp.CommandOutput, expectedLine) {
		t.Errorf("Expected environment output to contain %q, but it didn't.", expectedLine)
	}
}
