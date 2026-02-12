package create_file

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"net/http"
	"testing"
)

func TestCreateFile_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/create_file_success.txt"
	content := "hello"
	req := create_models.Request{
		Path:    filePath,
		Content: content,
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.CreateFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify using ReadFile API since filesystems are isolated
	readReq := read_models.Request{Path: filePath}
	resp, err := client.ReadFile(readReq)
	if err != nil {
		t.Fatalf("Failed to verify file via API: %v", err)
	}
	if resp.Content != content {
		t.Errorf("Expected content %q, got %q", content, resp.Content)
	}
}

func TestCreateFile_Conflict(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/create_file_conflict.txt"
	req := create_models.Request{
		Path:    filePath,
		Content: "already here",
	}

	// Create initial file via API
	_, err := client.CreateFile(req)
	if err != nil {
		t.Fatalf("Failed to setup initial file: %v", err)
	}

	req.Content = "new content"

	// -------------------------------------- Act --------------------------------------
	_, err = client.CreateFile(req)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusConflict, "File already exists")
}
