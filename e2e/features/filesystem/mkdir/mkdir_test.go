package mkdir

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	mkdir_models "agent-dev-environment/src/api/v1/filesystem/mkdir"
	"net/http"
	"testing"
)

func TestMkdir_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	dirPath := TestDir + "/mkdir_success/nested/dir"
	req := mkdir_models.Request{Path: dirPath}

	// -------------------------------------- Act --------------------------------------
	_, err := client.Mkdir(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify by creating a file inside it via API
	fileReq := create_models.Request{
		Path:    dirPath + "/test.txt",
		Content: "test",
	}
	_, err = client.CreateFile(fileReq)
	if err != nil {
		t.Fatalf("Failed to verify directory via API: %v", err)
	}
}

func TestMkdir_Idempotent(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	dirPath := TestDir + "/mkdir_idempotent"
	req := mkdir_models.Request{Path: dirPath}

	// Create it once
	_, err := client.Mkdir(req)
	if err != nil {
		t.Fatalf("First mkdir failed: %v", err)
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.Mkdir(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Second mkdir failed (should be idempotent): %v", err)
	}
}

func TestMkdir_ConflictWithFile(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/mkdir_conflict_file.txt"
	fileReq := create_models.Request{
		Path:    filePath,
		Content: "test",
	}
	_, err := client.CreateFile(fileReq)
	if err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	req := mkdir_models.Request{Path: filePath}

	// -------------------------------------- Act --------------------------------------
	_, err = client.Mkdir(req)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusConflict, "Path already exists and is not a directory")
}

func TestMkdir_BadRequest(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := mkdir_models.Request{Path: ""}

	// -------------------------------------- Act --------------------------------------
	_, err := client.Mkdir(req)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusBadRequest, "Path is required")
}
