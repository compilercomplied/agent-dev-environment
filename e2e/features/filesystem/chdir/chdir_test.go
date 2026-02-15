package chdir

import (
	"path/filepath"
	"strings"
	"testing"

	"agent-dev-environment/e2e"
	chdir_models "agent-dev-environment/src/api/v1/filesystem/chdir"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	ls_models "agent-dev-environment/src/api/v1/filesystem/ls"
	mkdir_models "agent-dev-environment/src/api/v1/filesystem/mkdir"
)

func TestChdir_ChangesWorkingDirectory(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	initialWdResp, err := client.Getwd()
	if err != nil {
		t.Fatalf("failed to get initial working directory: %v", err)
	}

	testDir := filepath.Join(e2e.TestDir, "chdir_basic_test")
	// Cleanup any leftovers before starting
	client.DeleteFile(delete_models.Request{Path: testDir, Recursive: true})

	defer func() {
		client.Chdir(chdir_models.Request{Path: initialWdResp.Path})
		client.DeleteFile(delete_models.Request{Path: testDir, Recursive: true})
	}()

	_, err = client.Mkdir(mkdir_models.Request{Path: testDir})
	if err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	markerFile := "marker.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, markerFile),
		Content: "here",
	})
	if err != nil {
		t.Fatalf("failed to create marker file: %v", err)
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.Chdir(chdir_models.Request{Path: testDir})

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

	if !strings.Contains(lsRes.CommandOutput, markerFile) {
		t.Errorf("expected output to contain %q, got: %s", markerFile, lsRes.CommandOutput)
	}
}

func TestChdir_AffectsSubsequentFileCreation(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	initialWdResp, err := client.Getwd()
	if err != nil {
		t.Fatalf("failed to get initial working directory: %v", err)
	}

	testDir := filepath.Join(e2e.TestDir, "chdir_create_test")
	client.DeleteFile(delete_models.Request{Path: testDir, Recursive: true})

	defer func() {
		client.Chdir(chdir_models.Request{Path: initialWdResp.Path})
		client.DeleteFile(delete_models.Request{Path: testDir, Recursive: true})
	}()

	_, _ = client.Mkdir(mkdir_models.Request{Path: testDir})
	_, _ = client.Chdir(chdir_models.Request{Path: testDir})
	
	fileName := "relative-created.txt"

	// -------------------------------------- Act --------------------------------------
	_, err = client.CreateFile(create_models.Request{
		Path:    fileName,
		Content: "created via chdir",
	})

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify file existence via API ListFiles
	lsRes, err := client.ListFiles(ls_models.Request{
		Path: ".",
	})
	if err != nil {
		t.Fatalf("failed to list current directory: %v", err)
	}

	if !strings.Contains(lsRes.CommandOutput, fileName) {
		t.Errorf("expected file %q to be listed in current directory, but it wasn't", fileName)
	}
}

func TestChdir_PathNotFound(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	nonExistentPath := filepath.Join(e2e.TestDir, "nonexistent_dir_for_chdir")

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
