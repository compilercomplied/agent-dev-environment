package reload_env

import (
	. "agent-dev-environment/e2e"
	replace_models "agent-dev-environment/src/api/v1/filesystem/replace"
	run_models "agent-dev-environment/src/api/v1/shell/run"
	"strings"
	"testing"
)

func TestReloadEnv_UpdatesProcessEnvironment(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	scriptPath := "scripts/load-env.sh"
	testVarName := "E2E_RELOAD_SUCCESS"
	testVarValue := "confirmed"
	// Injection line to be added to the script
	injection := "\necho \"" + testVarName + "='" + testVarValue + "'\" >> .env\n"
	oldString := "echo \"Environment loaded successfully from local stack\""
	newString := injection + oldString

	_, _ = client.ReloadEnv()

	_, err := client.Replace(replace_models.Request{
		Path:      scriptPath,
		OldString: oldString,
		NewString: newString,
	})
	if err != nil {
		t.Fatalf("Failed to arrange: could not modify script: %v", err)
	}

	defer func() {
		_, _ = client.Replace(replace_models.Request{
			Path:      scriptPath,
			OldString: newString,
			NewString: oldString,
		})
	}()

	// -------------------------------------- Act --------------------------------------
	_, err = client.ReloadEnv()

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Act failed: ReloadEnv returned error: %v", err)
	}

	resp, err := client.RunShell(run_models.Request{
		Command: "env",
		Args:    []string{},
	})

	if err != nil {
		t.Fatalf("Failed to assert: shell command failed: %v", err)
	}

	expectedLine := testVarName + "=" + testVarValue
	if !strings.Contains(resp.CommandOutput, expectedLine) {
		t.Errorf("Expected environment output to contain %q, but it didn't. Full output:\n%s", expectedLine, resp.CommandOutput)
	}
}
