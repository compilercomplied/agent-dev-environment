package replace

import (
	. "agent-dev-environment/e2e"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	replace_models "agent-dev-environment/src/api/v1/filesystem/replace"
	"net/http"
	"testing"
)

func TestReplace_ExactMatch_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/replace_exact_success.txt"
	initialContent := `Hello world!
This is a test.`
	oldString := "world"
	newString := "Gemini"
	expectedContent := `Hello Gemini!
This is a test.`

	_, err := client.CreateFile(create_models.Request{
		Path:    filePath,
		Content: initialContent,
	})
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}

	req := replace_models.Request{
		Path:      filePath,
		OldString: oldString,
		NewString: newString,
	}

	// -------------------------------------- Act --------------------------------------
	resp, err := client.Replace(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	readResp, err := client.ReadFile(read_models.Request{Path: filePath})
	if err != nil {
		t.Fatalf("Failed to read file after replacement: %v", err)
	}
	if readResp.Content != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, readResp.Content)
	}
}

func TestReplace_FuzzyMatch_Success(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/replace_fuzzy_success.txt"
	// 100+ characters long string
	initialContent := "The quick brown fox jumps over the lazy dog. This is a long string to test 98 percent similarity matches."
	// Provided string has one character different ('.' -> '!')
	providedOldString := "The quick brown fox jumps over the lazy dog! This is a long string to test 98 percent similarity matches."
	newString := "A fast brown fox"
	
	_, err := client.CreateFile(create_models.Request{
		Path:    filePath,
		Content: initialContent,
	})
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}

	req := replace_models.Request{
		Path:      filePath,
		OldString: providedOldString,
		NewString: newString,
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.Replace(req)

	// ------------------------------------ Assert -------------------------------------
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	readResp, err := client.ReadFile(read_models.Request{Path: filePath})
	if err != nil {
		t.Fatalf("Failed to read file after replacement: %v", err)
	}
	expectedContent := "A fast brown fox"
	if readResp.Content != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, readResp.Content)
	}
}

func TestReplace_FuzzyMatch_Failure(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	filePath := TestDir + "/replace_fuzzy_failure.txt"
	initialContent := "The quick brown fox jumps over the lazy dog."
	// String is 44 characters long.
	// 2 characters difference would be (44-2)/44 = 42/44 = 0.954... < 0.98
	providedOldString := "The quick brown fox jumps OVER the lazy dog!" 
	
	_, err := client.CreateFile(create_models.Request{
		Path:    filePath,
		Content: initialContent,
	})
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}

	req := replace_models.Request{
		Path:      filePath,
		OldString: providedOldString,
		NewString: "Something else",
	}

	// -------------------------------------- Act --------------------------------------
	_, err = client.Replace(req)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusBadRequest, "Could not find a match with at least 98% similarity")
}

func TestReplace_NotFound(t *testing.T) {
	// ------------------------------------ Arrange ------------------------------------
	client := NewClient()
	req := replace_models.Request{
		Path:      TestDir + "/non-existent.txt",
		OldString: "foo",
		NewString: "bar",
	}

	// -------------------------------------- Act --------------------------------------
	_, err := client.Replace(req)

	// ------------------------------------ Assert -------------------------------------
	AssertError(t, err, http.StatusNotFound, "File not found")
}
