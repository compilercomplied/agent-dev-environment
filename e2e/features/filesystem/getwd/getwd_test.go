package getwd

import (
	"path/filepath"
	"strings"
	"testing"

	"agent-dev-environment/e2e"
	chdir_models "agent-dev-environment/src/api/v1/filesystem/chdir"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	mkdir_models "agent-dev-environment/src/api/v1/filesystem/mkdir"
)

func TestGetwd_ReturnsCurrentDirectory(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	// -------------------------------------- Act --------------------------------------
	resp, err := client.Getwd()

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Path == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetwd_ReflectsChdir(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	initialWdResp, err := client.Getwd()
	if err != nil {
		t.Fatalf("failed to get initial working directory: %v", err)
	}

	testDirName := "getwd_reflects_chdir_test"
	testDir := filepath.Join(e2e.TestDir, testDirName)
	client.DeleteFile(delete_models.Request{Path: testDir, Recursive: true})

	defer func() {
		client.Chdir(chdir_models.Request{Path: initialWdResp.Path})
		client.DeleteFile(delete_models.Request{Path: testDir, Recursive: true})
	}()

	_, err = client.Mkdir(mkdir_models.Request{Path: testDir})
	if err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	_, err = client.Chdir(chdir_models.Request{Path: testDir})
	if err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.Getwd()

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.HasSuffix(resp.Path, testDirName) {
		t.Errorf("expected path to end with %q, got %q", testDirName, resp.Path)
	}
}
