package move_test

import (
	"net/http"
	"os"
	"testing"

	"agent-dev-environment/e2e"
	move_models "agent-dev-environment/src/api/v1/filesystem/move"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
)

func TestMoveFile_Success(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/test_move_source.txt"
	destPath := e2e.TestDir + "/test_move_dest.txt"
	content := "test content for move"

	// Create source file
	_, err := client.CreateFile(create_models.Request{
		Path:    sourcePath,
		Content: content,
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(destPath)

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// ---- Act ----
	_, err = client.MoveFile(req)

	// ---- Assert ----
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
		t.Errorf("Expected source file to be removed, but it still exists")
	}

	// Verify destination exists with correct content
	destContent, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(destContent) != content {
		t.Errorf("Expected destination content to be %q, got %q", content, string(destContent))
	}
}

func TestMoveFile_SourceNotFound(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/nonexistent_source.txt"
	destPath := e2e.TestDir + "/test_move_dest2.txt"

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// ---- Act ----
	_, err := client.MoveFile(req)

	// ---- Assert ----
	e2e.AssertError(t, err, http.StatusNotFound, "Source path does not exist")
}

func TestMoveFile_DestinationAlreadyExists(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/test_move_source3.txt"
	destPath := e2e.TestDir + "/test_move_dest3.txt"

	// Create both source and destination files
	_, err := client.CreateFile(create_models.Request{
		Path:    sourcePath,
		Content: "source content",
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(sourcePath)

	_, err = client.CreateFile(create_models.Request{
		Path:    destPath,
		Content: "destination content",
	})
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer os.Remove(destPath)

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// ---- Act ----
	_, err = client.MoveFile(req)

	// ---- Assert ----
	e2e.AssertError(t, err, http.StatusConflict, "Destination path already exists")
}

func TestMoveFile_EmptySource(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	destPath := e2e.TestDir + "/test_move_dest4.txt"

	req := move_models.Request{
		Source:      "",
		Destination: destPath,
	}

	// ---- Act ----
	_, err := client.MoveFile(req)

	// ---- Assert ----
	e2e.AssertError(t, err, http.StatusBadRequest, "Source path is required")
}

func TestMoveFile_EmptyDestination(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/test_move_source5.txt"

	req := move_models.Request{
		Source:      sourcePath,
		Destination: "",
	}

	// ---- Act ----
	_, err := client.MoveFile(req)

	// ---- Assert ----
	e2e.AssertError(t, err, http.StatusBadRequest, "Destination path is required")
}

func TestMoveDirectory_Success(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourceDir := e2e.TestDir + "/test_move_dir_source"
	destDir := e2e.TestDir + "/test_move_dir_dest"

	// Create source directory with a file
	err := os.Mkdir(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	testFilePath := sourceDir + "/testfile.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    testFilePath,
		Content: "test content in directory",
	})
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.RemoveAll(destDir)

	req := move_models.Request{
		Source:      sourceDir,
		Destination: destDir,
	}

	// ---- Act ----
	_, err = client.MoveFile(req)

	// ---- Assert ----
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	if _, err := os.Stat(sourceDir); !os.IsNotExist(err) {
		t.Errorf("Expected source directory to be removed, but it still exists")
	}

	// Verify destination exists and contains the file
	destFilePath := destDir + "/testfile.txt"
	destContent, err := os.ReadFile(destFilePath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(destContent) != "test content in directory" {
		t.Errorf("Expected destination file content to be %q, got %q", "test content in directory", string(destContent))
	}
}
