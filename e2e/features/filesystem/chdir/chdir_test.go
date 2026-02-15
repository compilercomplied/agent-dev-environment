package chdir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"agent-dev-environment/e2e"
	chdir_models "agent-dev-environment/src/api/v1/filesystem/chdir"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	ls_models "agent-dev-environment/src/api/v1/filesystem/ls"
)

func TestChdir_ChangesWorkingDirectory(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	// Capture initial working directory to restore it later
	initialWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() {
		_, _ = client.Chdir(chdir_models.Request{Path: initialWd})
	}()

	// Create a temporary directory for testing chdir
	tmpDir, err := os.MkdirTemp("", "chdir-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file in the temp dir via OS to verify chdir later
	fileName := "marker.txt"
	err = os.WriteFile(filepath.Join(tmpDir, fileName), []byte("here"), 0644)
	if err != nil {
		t.Fatalf("failed to create marker file: %v", err)
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.Chdir(chdir_models.Request{Path: tmpDir})

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the change by listing "."
	lsRes, err := client.ListFiles(ls_models.Request{
		Path: ".",
	})
	if err != nil {
		t.Fatalf("failed to list current directory: %v", err)
	}

	if !strings.Contains(lsRes.CommandOutput, fileName) {
		t.Errorf("expected output to contain %q, got: %s", fileName, lsRes.CommandOutput)
	}
}

func TestChdir_AffectsSubsequentFileCreation(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	initialWd, _ := os.Getwd()
	defer func() {
		_, _ = client.Chdir(chdir_models.Request{Path: initialWd})
	}()

	tmpDir, _ := os.MkdirTemp("", "chdir-create-test-*")
	defer os.RemoveAll(tmpDir)

	_, _ = client.Chdir(chdir_models.Request{Path: tmpDir})
	
	fileName := "relative-created.txt"

	// -------------------------------------- Act --------------------------------------
	_, err := client.CreateFile(create_models.Request{
		Path:    fileName,
		Content: "created via chdir",
	})

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify file existence in the correct physical directory
	expectedPath := filepath.Join(tmpDir, fileName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s, but it doesn't", expectedPath)
	}
}

func TestChdir_PathNotFound(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	nonExistentPath := "/tmp/path/that/does/not/exist/anywhere"

	// -------------------------------------- Act --------------------------------------
	_, err := client.Chdir(chdir_models.Request{Path: nonExistentPath})

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, 404, "Directory not found")
}

func TestChdir_EmptyPath(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	// -------------------------------------- Act --------------------------------------
	_, err := client.Chdir(chdir_models.Request{Path: ""})

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, 400, "Path is required")
}
