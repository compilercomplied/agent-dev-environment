package read

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"strings"
	"testing"
)

func TestRead_Pagination(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/pagination_test.txt"
	lines := []string{"line1", "line2", "line3", "line4", "line5"}
	content := strings.Join(lines, "\n")
	
	_, err := client.CreateFile(create_models.Request{
		Path:    filePath,
		Content: content,
	})
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("should read first 2 lines", func(t *testing.T) {
		// ------------------------------------ Arrange ------------------------------------
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
		if resp.LinesRead != 2 {
			t.Errorf("Expected LinesRead 2, got %d", resp.LinesRead)
		}
	})

	t.Run("should read from middle to end", func(t *testing.T) {
		// ------------------------------------ Arrange ------------------------------------
		offset := 2
		limit := 10 // More than remaining
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
		if resp.LinesRead != 3 {
			t.Errorf("Expected LinesRead 3, got %d", resp.LinesRead)
		}
	})

	t.Run("should read whole file when no pagination provided", func(t *testing.T) {
		// ------------------------------------ Arrange ------------------------------------
		req := read_models.Request{Path: filePath}

		// -------------------------------------- Act --------------------------------------
		resp, err := client.ReadFile(req)

		// ------------------------------------ Assert -------------------------------------
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.Content != content {
			t.Errorf("Expected full content, got %q", resp.Content)
		}
		if resp.TotalLines != 5 {
			t.Errorf("Expected TotalLines 5, got %d", resp.TotalLines)
		}
		if resp.HasMore {
			t.Error("Expected HasMore to be false")
		}
	})
}
