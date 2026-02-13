package move_test

import (
	"net/http"
	"testing"

	"agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	move_models "agent-dev-environment/src/api/v1/filesystem/move"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
)

func TestMoveFile_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
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
	defer client.DeleteFile(delete_models.Request{Path: destPath})

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	_, err = client.ReadFile(read_models.Request{Path: sourcePath})
	e2e.AssertError(t, err, http.StatusNotFound, "File not found")

	// Verify destination exists with correct content
	resp, err := client.ReadFile(read_models.Request{Path: destPath})
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if resp.Content != content {
		t.Errorf("Expected destination content to be %q, got %q", content, resp.Content)
	}
}

func TestMoveFile_SourceNotFound(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/nonexistent_source.txt"
	destPath := e2e.TestDir + "/test_move_dest2.txt"

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, http.StatusNotFound, "Source path does not exist")
}

func TestMoveFile_DestinationAlreadyExists(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
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
	defer client.DeleteFile(delete_models.Request{Path: sourcePath})

	_, err = client.CreateFile(create_models.Request{
		Path:    destPath,
		Content: "destination content",
	})
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer client.DeleteFile(delete_models.Request{Path: destPath})

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, http.StatusConflict, "Destination path already exists")
}

func TestMoveFile_EmptySource(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	destPath := e2e.TestDir + "/test_move_dest4.txt"

	req := move_models.Request{
		Source:      "",
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, http.StatusBadRequest, "Source path is required")
}

func TestMoveFile_EmptyDestination(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	sourcePath := e2e.TestDir + "/test_move_source5.txt"

	req := move_models.Request{
		Source:      sourcePath,
		Destination: "",
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, http.StatusBadRequest, "Destination path is required")
}

func TestMoveDirectory_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	sourceDir := e2e.TestDir + "/test_move_dir_source"
	destDir := e2e.TestDir + "/test_move_dir_dest"

	// Create source directory with a file via API (auto-creates parent dirs)
	testFilePath := sourceDir + "/testfile.txt"
	_, err := client.CreateFile(create_models.Request{
		Path:    testFilePath,
		Content: "test content in directory",
	})
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer client.DeleteFile(delete_models.Request{Path: destDir, Recursive: true})

	req := move_models.Request{
		Source:      sourceDir,
		Destination: destDir,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	_, err = client.ReadFile(read_models.Request{Path: testFilePath})
	e2e.AssertError(t, err, http.StatusNotFound, "File not found")

	// Verify destination exists and contains the file
	destFilePath := destDir + "/testfile.txt"
	resp, err := client.ReadFile(read_models.Request{Path: destFilePath})
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if resp.Content != "test content in directory" {
		t.Errorf("Expected destination file content to be %q, got %q", "test content in directory", resp.Content)
	}
}

func TestMoveFile_RenameInSameDirectory(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
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
	defer client.DeleteFile(delete_models.Request{Path: destPath})

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	_, err = client.ReadFile(read_models.Request{Path: sourcePath})
	e2e.AssertError(t, err, http.StatusNotFound, "File not found")

	// Verify destination exists with correct content
	resp, err := client.ReadFile(read_models.Request{Path: destPath})
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if resp.Content != content {
		t.Errorf("Expected destination content to be %q, got %q", content, resp.Content)
	}
}

func TestMoveFile_AcrossNestedDirectories(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	sourceDir := e2e.TestDir + "/nested_source"
	destDir := e2e.TestDir + "/nested_dest"

	sourcePath := sourceDir + "/file.txt"
	destPath := destDir + "/file.txt"
	content := "nested file content"

	// Create source file via API (auto-creates parent dirs)
	_, err := client.CreateFile(create_models.Request{
		Path:    sourcePath,
		Content: content,
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer client.DeleteFile(delete_models.Request{Path: sourceDir, Recursive: true})

	// Create destination directory via API by creating and deleting a temp file
	destPlaceholder := destDir + "/.placeholder"
	_, err = client.CreateFile(create_models.Request{
		Path:    destPlaceholder,
		Content: "",
	})
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}
	client.DeleteFile(delete_models.Request{Path: destPlaceholder})
	defer client.DeleteFile(delete_models.Request{Path: destDir, Recursive: true})

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	_, err = client.ReadFile(read_models.Request{Path: sourcePath})
	e2e.AssertError(t, err, http.StatusNotFound, "File not found")

	// Verify destination exists with correct content
	resp, err := client.ReadFile(read_models.Request{Path: destPath})
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if resp.Content != content {
		t.Errorf("Expected destination content to be %q, got %q", content, resp.Content)
	}
}

func TestMoveDirectory_WithMultipleFiles(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	sourceDir := e2e.TestDir + "/multi_file_source"
	destDir := e2e.TestDir + "/multi_file_dest"

	// Create multiple files via API (auto-creates parent dirs including subdir)
	file1Path := sourceDir + "/file1.txt"
	_, err := client.CreateFile(create_models.Request{
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

	file3Path := sourceDir + "/subdir/file3.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    file3Path,
		Content: "content 3",
	})
	if err != nil {
		t.Fatalf("Failed to create file3: %v", err)
	}

	defer client.DeleteFile(delete_models.Request{Path: destDir, Recursive: true})

	req := move_models.Request{
		Source:      sourceDir,
		Destination: destDir,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify source no longer exists
	_, err = client.ReadFile(read_models.Request{Path: file1Path})
	e2e.AssertError(t, err, http.StatusNotFound, "File not found")

	// Verify all files exist in destination
	resp1, err := client.ReadFile(read_models.Request{Path: destDir + "/file1.txt"})
	if err != nil {
		t.Errorf("Failed to read destination file1: %v", err)
	} else if resp1.Content != "content 1" {
		t.Errorf("Expected file1 content to be %q, got %q", "content 1", resp1.Content)
	}

	resp2, err := client.ReadFile(read_models.Request{Path: destDir + "/file2.txt"})
	if err != nil {
		t.Errorf("Failed to read destination file2: %v", err)
	} else if resp2.Content != "content 2" {
		t.Errorf("Expected file2 content to be %q, got %q", "content 2", resp2.Content)
	}

	resp3, err := client.ReadFile(read_models.Request{Path: destDir + "/subdir/file3.txt"})
	if err != nil {
		t.Errorf("Failed to read destination file3: %v", err)
	} else if resp3.Content != "content 3" {
		t.Errorf("Expected file3 content to be %q, got %q", "content 3", resp3.Content)
	}
}

func TestMoveFile_DestinationParentDoesNotExist(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
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
	defer client.DeleteFile(delete_models.Request{Path: sourcePath})

	req := move_models.Request{
		Source:      sourcePath,
		Destination: destPath,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err == nil {
		t.Fatal("Expected error when destination parent directory doesn't exist, got nil")
	}

	// Verify source still exists (move should have failed)
	_, readErr := client.ReadFile(read_models.Request{Path: sourcePath})
	if readErr != nil {
		t.Errorf("Expected source file to still exist after failed move, got: %v", readErr)
	}
}

func TestMoveDirectory_DestinationAlreadyExists(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()

	sourceDir := e2e.TestDir + "/dir_exists_source"
	destDir := e2e.TestDir + "/dir_exists_dest"

	// Create source directory with a file via API
	sourceFile := sourceDir + "/file.txt"
	_, err := client.CreateFile(create_models.Request{
		Path:    sourceFile,
		Content: "source content",
	})
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer client.DeleteFile(delete_models.Request{Path: sourceDir, Recursive: true})

	// Create destination directory with a file via API
	destFile := destDir + "/existing.txt"
	_, err = client.CreateFile(create_models.Request{
		Path:    destFile,
		Content: "dest content",
	})
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer client.DeleteFile(delete_models.Request{Path: destDir, Recursive: true})

	req := move_models.Request{
		Source:      sourceDir,
		Destination: destDir,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.MoveFile(req)

	// ------------------------------------ Assert -------------------------------------
	e2e.AssertError(t, err, http.StatusConflict, "Destination path already exists")
}
