package ls

import (
	"path/filepath"
	"testing"

	"agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	ls_models "agent-dev-environment/src/api/v1/filesystem/ls"
)

func TestListFiles_BasicDirectory(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_basic")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file1.txt"),
		Content: "content1",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file2.txt"),
		Content: "content2",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file3.txt"),
		Content: "content3",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(resp.Entries))
	}

	expectedFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for i, entry := range resp.Entries {
		if entry.Name != expectedFiles[i] {
			t.Errorf("expected file %s, got %s", expectedFiles[i], entry.Name)
		}
		if entry.IsDirectory {
			t.Errorf("expected file %s to not be a directory", entry.Name)
		}
	}
}

func TestListFiles_WithLongFormat(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_long")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file.txt"),
		Content: "test content",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      true,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(resp.Entries))
	}

	entry := resp.Entries[0]
	if entry.Name != "file.txt" {
		t.Errorf("expected file name 'file.txt', got %s", entry.Name)
	}
	if entry.Size != 12 {
		t.Errorf("expected size 12, got %d", entry.Size)
	}
	if entry.Mode == "" {
		t.Error("expected mode to be set in long format")
	}
	if entry.ModTime.IsZero() {
		t.Error("expected mod_time to be set in long format")
	}
}

func TestListFiles_WithoutLongFormat(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_no_long")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file.txt"),
		Content: "test content",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(resp.Entries))
	}

	entry := resp.Entries[0]
	if entry.Size != 12 {
		t.Errorf("expected size 12 to always be present, got %d", entry.Size)
	}
	if entry.Mode != "" {
		t.Error("expected mode to be empty without long format")
	}
	if !entry.ModTime.IsZero() {
		t.Error("expected mod_time to be zero without long format")
	}
}

func TestListFiles_Recursive(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_recursive")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file1.txt"),
		Content: "content1",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "subdir", "file2.txt"),
		Content: "content2",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "subdir", "nested", "file3.txt"),
		Content: "content3",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: true,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) < 3 {
		t.Fatalf("expected at least 3 file entries, got %d", len(resp.Entries))
	}

	fileNames := make(map[string]bool)
	for _, entry := range resp.Entries {
		fileNames[entry.Name] = true
	}

	expectedFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, fileName := range expectedFiles {
		if !fileNames[fileName] {
			t.Errorf("expected to find file %s in recursive listing", fileName)
		}
	}
}

func TestListFiles_NonRecursive(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_non_recursive")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file1.txt"),
		Content: "content1",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "subdir", "file2.txt"),
		Content: "content2",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) != 2 {
		t.Fatalf("expected 2 entries (1 file, 1 directory), got %d", len(resp.Entries))
	}

	fileNames := make(map[string]bool)
	for _, entry := range resp.Entries {
		fileNames[entry.Name] = true
		if entry.Name == "file2.txt" {
			t.Errorf("file2.txt should not be in non-recursive listing of parent directory")
		}
	}

	if !fileNames["file1.txt"] {
		t.Error("expected to find file1.txt in non-recursive listing")
	}
	if !fileNames["subdir"] {
		t.Error("expected to find subdir in non-recursive listing")
	}
}

func TestListFiles_PathNotFound(t *testing.T) {
	client := e2e.NewClient()
	testPath := filepath.Join(e2e.TestDir, "nonexistent_dir")

	_, err := client.ListFiles(ls_models.Request{
		Path:      testPath,
		Recursive: false,
		Long:      false,
	})

	e2e.AssertError(t, err, 404, "Path not found")
}

func TestListFiles_EmptyPath(t *testing.T) {
	client := e2e.NewClient()

	_, err := client.ListFiles(ls_models.Request{
		Path:      "",
		Recursive: false,
		Long:      false,
	})

	e2e.AssertError(t, err, 400, "Path is required")
}

func TestListFiles_SingleFile(t *testing.T) {
	client := e2e.NewClient()
	testFile := filepath.Join(e2e.TestDir, "test_ls_single_file.txt")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testFile,
			Recursive: false,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    testFile,
		Content: "test content",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testFile,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(resp.Entries))
	}

	entry := resp.Entries[0]
	if entry.Name != "test_ls_single_file.txt" {
		t.Errorf("expected file name 'test_ls_single_file.txt', got %s", entry.Name)
	}
	if entry.IsDirectory {
		t.Error("expected entry to not be a directory")
	}
}

func TestListFiles_EmptyDirectory(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_empty_dir")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, ".placeholder"),
		Content: "",
	})
	client.DeleteFile(delete_models.Request{
		Path:      filepath.Join(testDir, ".placeholder"),
		Recursive: false,
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Entries) != 0 {
		t.Fatalf("expected 0 entries in empty directory, got %d", len(resp.Entries))
	}
}

func TestListFiles_IncludesDirectories(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_with_dirs")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file.txt"),
		Content: "content",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "subdir1", ".placeholder"),
		Content: "",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "subdir2", ".placeholder"),
		Content: "",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var fileCount, dirCount int
	for _, entry := range resp.Entries {
		if entry.IsDirectory {
			dirCount++
		} else {
			fileCount++
		}
	}

	if fileCount != 1 {
		t.Errorf("expected 1 file, got %d", fileCount)
	}
	if dirCount != 2 {
		t.Errorf("expected 2 directories, got %d", dirCount)
	}
}

func TestListFiles_SizeAlwaysPresent(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_size")

	client.DeleteFile(delete_models.Request{
		Path:      testDir,
		Recursive: true,
	})

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "small.txt"),
		Content: "a",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "large.txt"),
		Content: "this is a larger file with more content",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, entry := range resp.Entries {
		if entry.Name == "small.txt" && entry.Size != 1 {
			t.Errorf("expected small.txt size to be 1, got %d", entry.Size)
		}
		if entry.Name == "large.txt" && entry.Size != 39 {
			t.Errorf("expected large.txt size to be 39, got %d", entry.Size)
		}
	}
}
