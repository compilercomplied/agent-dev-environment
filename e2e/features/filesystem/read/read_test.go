package read

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"net/http"
	"testing"
)

func TestReadFile_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/read_file_success.txt"
	content := "some content to read"

	// Create file via API
	_, err := client.CreateFile(create_models.Request{
		Path:    filePath,
		Content: content,
	})
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}

	req := read_models.Request{
		Path: filePath,
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.Content != content {
		t.Errorf("Expected content %q, got %q", content, resp.Content)
	}
}

func TestReadFile_NotFound(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := read_models.Request{
		Path: TestDir + "/this-file-really-should-not-exist-12345.txt",
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusNotFound, "File not found")
}
