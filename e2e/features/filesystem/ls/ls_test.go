package ls

import (
	"path/filepath"
	"strings"
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

	got := strings.TrimSpace(resp.CommandOutput)
	expected := "file1.txt\nfile2.txt\nfile3.txt"
	if got != expected {
		t.Errorf("expected exactly %q, got %q", expected, got)
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

	output := resp.CommandOutput
	if output == "" {
		t.Fatal("expected command output to be non-empty")
	}

	// Long format should contain file permissions
	if !strings.Contains(output, "file.txt") {
		t.Errorf("expected output to contain 'file.txt', got: %s", output)
	}

	// Long format should have permission info (starting with - or d)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	foundFileInfo := false
	for _, line := range lines {
		if strings.Contains(line, "file.txt") && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "d")) {
			foundFileInfo = true
			break
		}
	}
	if !foundFileInfo {
		t.Errorf("expected long format output with permissions, got: %s", output)
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

	output := resp.CommandOutput
	if output == "" {
		t.Fatal("expected command output to be non-empty")
	}

	// Should contain file name
	if !strings.Contains(output, "file.txt") {
		t.Errorf("expected output to contain 'file.txt', got: %s", output)
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

	output := resp.CommandOutput
	if output == "" {
		t.Fatal("expected command output to be non-empty")
	}

	// All files should be present in recursive output
	expectedFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range expectedFiles {
		if !strings.Contains(output, file) {
			t.Errorf("expected recursive output to contain %s, got: %s", file, output)
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

	output := resp.CommandOutput
	if output == "" {
		t.Fatal("expected command output to be non-empty")
	}

	// Should contain file1.txt and subdir
	if !strings.Contains(output, "file1.txt") {
		t.Errorf("expected output to contain 'file1.txt', got: %s", output)
	}
	if !strings.Contains(output, "subdir") {
		t.Errorf("expected output to contain 'subdir', got: %s", output)
	}

	// Should NOT contain file2.txt (it's in subdir)
	if strings.Contains(output, "file2.txt") {
		t.Errorf("expected non-recursive output NOT to contain 'file2.txt', got: %s", output)
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

	output := resp.CommandOutput
	if output == "" {
		t.Fatal("expected command output to be non-empty")
	}

	// When listing a single file, ls returns the file path
	if !strings.Contains(output, testFile) {
		t.Errorf("expected output to contain the file path, got: %s", output)
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

	// Empty directory output should be empty or just whitespace
	output := strings.TrimSpace(resp.CommandOutput)
	if output != "" {
		t.Fatalf("expected empty output for empty directory, got: %s", output)
	}
}

func TestListFiles_CommandOutputField(t *testing.T) {
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_ls_output_field")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "test.txt"),
		Content: "content",
	})

	resp, err := client.ListFiles(ls_models.Request{
		Path:      testDir,
		Recursive: false,
		Long:      false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the response has the CommandOutput field with raw ls output
	if resp.CommandOutput == "" {
		t.Error("expected CommandOutput field to be populated with ls command output")
	}

	// The output should be a raw string from the ls command
	if !strings.Contains(resp.CommandOutput, "test.txt") {
		t.Errorf("expected CommandOutput to contain raw ls output with 'test.txt', got: %s", resp.CommandOutput)
	}
}
