package search

import (
	"path/filepath"
	"strings"
	"testing"

	"agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	search_models "agent-dev-environment/src/api/v1/filesystem/search"
)

func TestSearch_Basic(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_search_basic")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file1.txt"),
		Content: `Hello world
This is a test
`,
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file2.txt"),
		Content: `Another file
No match here
`,
	})

	// -------------------------------------- Act --------------------------------------
	resp, err := client.Search(search_models.Request{
		Path:    testDir,
		Pattern: "Hello",
	})

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := strings.TrimSpace(resp.CommandOutput)
	expected := filepath.Join(testDir, "file1.txt") + ":Hello world"
	if got != expected {
		t.Errorf("expected exactly %q, got %q", expected, got)
	}
}

func TestSearch_IgnoreCase(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_search_ignore_case")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "file1.txt"),
		Content: "Hello World\n",
	})

	// -------------------------------------- Act --------------------------------------
	// Case-sensitive search (default)
	resp, err := client.Search(search_models.Request{
		Path:    testDir,
		Pattern: "hello",
	})
	
	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strings.TrimSpace(resp.CommandOutput) != "" {
		t.Errorf("expected empty output for case-sensitive mismatch, got %q", resp.CommandOutput)
	}

	// -------------------------------------- Act --------------------------------------
	// Case-insensitive search
	resp, err = client.Search(search_models.Request{
		Path:       testDir,
		Pattern:    "hello",
		IgnoreCase: true,
	})
	
	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	got := strings.TrimSpace(resp.CommandOutput)
	expected := filepath.Join(testDir, "file1.txt") + ":Hello World"
	if got != expected {
		t.Errorf("expected exactly %q, got %q", expected, got)
	}
}

func TestSearch_FilesWithMatches(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := e2e.NewClient()
	testDir := filepath.Join(e2e.TestDir, "test_search_files_with_matches")

	defer func() {
		client.DeleteFile(delete_models.Request{
			Path:      testDir,
			Recursive: true,
		})
	}()

	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "match.txt"),
		Content: "Target found here\n",
	})
	client.CreateFile(create_models.Request{
		Path:    filepath.Join(testDir, "no_match.txt"),
		Content: "Nothing here\n",
	})

	// -------------------------------------- Act --------------------------------------
	resp, err := client.Search(search_models.Request{
		Path:             testDir,
		Pattern:          "Target",
		FilesWithMatches: true,
	})

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := strings.TrimSpace(resp.CommandOutput)
	expected := filepath.Join(testDir, "match.txt")
	if got != expected {
		t.Errorf("expected exactly %q, got %q", expected, got)
	}
}
