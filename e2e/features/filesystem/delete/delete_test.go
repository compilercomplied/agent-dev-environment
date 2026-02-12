package delete

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"net/http"
	"testing"
)

func TestDeleteFile_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/delete_file_success.txt"
	createReq := create_models.Request{
		Path:    filePath,
		Content: "test content",
	}

	// Create file via API
	_, err := client.CreateFile(createReq)
	if err != nil {
		t.Fatalf("Failed to setup initial file: %v", err)
	}

	deleteReq := delete_models.Request{
		Path:      filePath,
		Recursive: false,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.DeleteFile(deleteReq)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify file is deleted by trying to read it
	readReq := read_models.Request{Path: filePath}
	_, err = client.ReadFile(readReq)
	AssertError(t, err, http.StatusNotFound, "File not found")
}

func TestDeleteFile_NotFound(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/nonexistent.txt"
	deleteReq := delete_models.Request{
		Path:      filePath,
		Recursive: false,
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.DeleteFile(deleteReq)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusNotFound, "File or directory not found")
}

func TestDeleteDirectory_WithoutRecursiveFlag(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	dirPath := TestDir + "/test_dir_no_recursive"
	filePath := dirPath + "/file.txt"

	// Create directory and file via API
	createReq := create_models.Request{
		Path:    filePath,
		Content: "content",
	}
	_, err := client.CreateFile(createReq)
	if err != nil {
		t.Fatalf("Failed to setup directory and file: %v", err)
	}

	deleteReq := delete_models.Request{
		Path:      dirPath,
		Recursive: false,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.DeleteFile(deleteReq)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusBadRequest, "Cannot delete directory without recursive flag")
}

func TestDeleteDirectory_WithRecursiveFlag(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	dirPath := TestDir + "/test_dir_recursive"
	filePath := dirPath + "/file.txt"

	// Create directory and file via API
	createReq := create_models.Request{
		Path:    filePath,
		Content: "content",
	}
	_, err := client.CreateFile(createReq)
	if err != nil {
		t.Fatalf("Failed to setup directory and file: %v", err)
	}

	deleteReq := delete_models.Request{
		Path:      dirPath,
		Recursive: true,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.DeleteFile(deleteReq)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify directory and its contents are deleted
	readReq := read_models.Request{Path: filePath}
	_, err = client.ReadFile(readReq)
	AssertError(t, err, http.StatusNotFound, "File not found")
}

func TestDeleteFile_EmptyPath(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	deleteReq := delete_models.Request{
		Path:      "",
		Recursive: false,
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.DeleteFile(deleteReq)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusBadRequest, "Path is required")
}
