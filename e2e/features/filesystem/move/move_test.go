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

	// Ensure parent test directory exists
	err := os.MkdirAll(e2e.TestDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create source directory with a file
	err = os.Mkdir(sourceDir, 0755)
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

func TestMoveFile_RenameInSameDirectory(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/rename_source.txt"
	destPath := e2e.TestDir + "/rename_dest.txt"
	content := "content to rename"

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

func TestMoveFile_AcrossNestedDirectories(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()

	// Ensure parent test directory exists
	err := os.MkdirAll(e2e.TestDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	sourceDir := e2e.TestDir + "/nested_source"
	destDir := e2e.TestDir + "/nested_dest"

	err = os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}
	defer os.RemoveAll(destDir)

	sourcePath := sourceDir + "/file.txt"
	destPath := destDir + "/file.txt"
	content := "nested file content"

	// Create source file
	_, err = client.CreateFile(create_models.Request{
		Path:    sourcePath,
		Content: content,
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

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

func TestMoveDirectory_WithMultipleFiles(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()

	// Ensure parent test directory exists
	err := os.MkdirAll(e2e.TestDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	sourceDir := e2e.TestDir + "/multi_file_source"
	destDir := e2e.TestDir + "/multi_file_dest"

	// Create source directory
	err = os.Mkdir(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create subdirectory
	subDir := sourceDir + "/subdir"
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create multiple files
	file1Path := sourceDir + "/file1.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    file1Path,
		Content: "content 1",
	})
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	file2Path := sourceDir + "/file2.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    file2Path,
		Content: "content 2",
	})
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	file3Path := subDir + "/file3.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    file3Path,
		Content: "content 3",
	})
	if err != nil {
		t.Fatalf("Failed to create file3: %v", err)
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

	// Verify all files exist in destination
	destFile1 := destDir + "/file1.txt"
	content1, err := os.ReadFile(destFile1)
	if err != nil {
		t.Errorf("Failed to read destination file1: %v", err)
	} else if string(content1) != "content 1" {
		t.Errorf("Expected file1 content to be %q, got %q", "content 1", string(content1))
	}

	destFile2 := destDir + "/file2.txt"
	content2, err := os.ReadFile(destFile2)
	if err != nil {
		t.Errorf("Failed to read destination file2: %v", err)
	} else if string(content2) != "content 2" {
		t.Errorf("Expected file2 content to be %q, got %q", "content 2", string(content2))
	}

	destFile3 := destDir + "/subdir/file3.txt"
	content3, err := os.ReadFile(destFile3)
	if err != nil {
		t.Errorf("Failed to read destination file3: %v", err)
	} else if string(content3) != "content 3" {
		t.Errorf("Expected file3 content to be %q, got %q", "content 3", string(content3))
	}
}

func TestMoveFile_DestinationParentDoesNotExist(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/source_no_parent.txt"
	destPath := e2e.TestDir + "/nonexistent_dir/dest.txt"

	// Create source file
	_, err := client.CreateFile(create_models.Request{
		Path:    sourcePath,
		Content: "test content",
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(sourcePath)

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// ---- Act ----
	_, err = client.MoveFile(req)

	// ---- Assert ----
	if err == nil {
		t.Fatal("Expected error when destination parent directory doesn't exist, got nil")
	}

	// Verify source still exists (move should have failed)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		t.Errorf("Expected source file to still exist after failed move")
	}
}

func TestMoveDirectory_DestinationAlreadyExists(t *testing.T) {
	// ---- Arrange ----
	client := e2e.NewClient()

	// Ensure parent test directory exists
	err := os.MkdirAll(e2e.TestDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	sourceDir := e2e.TestDir + "/dir_exists_source"
	destDir := e2e.TestDir + "/dir_exists_dest"

	// Create source directory with a file
	err = os.Mkdir(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	sourceFile := sourceDir + "/file.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    sourceFile,
		Content: "source content",
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create destination directory
	err = os.Mkdir(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}
	defer os.RemoveAll(destDir)

	req := move_models.Request{
		Source:      sourceDir,
		Destination: destDir,
	}

	// ---- Act ----
	_, err = client.MoveFile(req)

	// ---- Assert ----
	e2e.AssertError(t, err, http.StatusConflict, "Destination path already exists")
}
