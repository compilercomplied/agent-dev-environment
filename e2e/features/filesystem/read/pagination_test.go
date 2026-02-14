package read

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"strings"
	"testing"
)

func setupPaginationTestFile(t *testing.T, client *Client, filePath string) string {
	lines := []string{"line1", "line2", "line3", "line4", "line5"}
	content := strings.Join(lines, "\n")
	_, err := client.CreateFile(create_models.Request{
		Path:    filePath,
		Content: content,
	})
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return content
}

func TestRead_Pagination_FirstLines(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/read_first_lines.txt"
	setupPaginationTestFile(t, client, filePath)
	offset := 0
	limit := 2
	req := read_models.Request{
		Path:   filePath,
		Offset: &offset,
		Limit:  &limit,
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected := "line1\nline2"
	if resp.Content != expected {
		t.Errorf("Expected content %q, got %q", expected, resp.Content)
	}
	if resp.TotalLines != 5 {
		t.Errorf("Expected TotalLines 5, got %d", resp.TotalLines)
	}
	if !resp.HasMore {
		t.Error("Expected HasMore to be true")
	}
}

func TestRead_Pagination_MiddleToEnd(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/read_middle_to_end.txt"
	setupPaginationTestFile(t, client, filePath)
	offset := 2
	limit := 10
	req := read_models.Request{
		Path:   filePath,
		Offset: &offset,
		Limit:  &limit,
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected := "line3\nline4\nline5"
	if resp.Content != expected {
		t.Errorf("Expected content %q, got %q", expected, resp.Content)
	}
	if resp.HasMore {
		t.Error("Expected HasMore to be false")
	}
}

func TestRead_Pagination_OutOfBoundsOffset(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/read_out_of_bounds.txt"
	setupPaginationTestFile(t, client, filePath)
	offset := 10
	limit := 5
	req := read_models.Request{
		Path:   filePath,
		Offset: &offset,
		Limit:  &limit,
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Content != "" {
		t.Errorf("Expected empty content for out of bounds offset, got %q", resp.Content)
	}
	if resp.LinesRead != 0 {
		t.Errorf("Expected 0 lines read, got %d", resp.LinesRead)
	}
}

func TestRead_Pagination_LimitOne(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/read_limit_one.txt"
	setupPaginationTestFile(t, client, filePath)
	offset := 4
	limit := 1
	req := read_models.Request{
		Path:   filePath,
		Offset: &offset,
		Limit:  &limit,
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Content != "line5" {
		t.Errorf("Expected 'line5', got %q", resp.Content)
	}
	if resp.HasMore {
		t.Error("Expected HasMore to be false for the last line")
	}
}

func TestRead_Pagination_FullFileNoParams(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/read_full_no_params.txt"
	content := setupPaginationTestFile(t, client, filePath)
	req := read_models.Request{Path: filePath}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.ReadFile(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Content != content {
		t.Error("Expected full content")
	}
	if resp.TotalLines != 5 {
		t.Errorf("Expected 5 total lines, got %d", resp.TotalLines)
	}
}
