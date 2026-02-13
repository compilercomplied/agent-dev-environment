package replace

import (
	"os"
	"strings"

	replace_models "agent-dev-environment/src/api/v1/filesystem/replace"
	"agent-dev-environment/src/library/api"
)

func Handler(req replace_models.Request) (*replace_models.Response, error) {
	content, err := os.ReadFile(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "File not found")
		}
		return nil, err
	}

	fileContent := string(content)

	// Try exact match first
	if strings.Contains(fileContent, req.OldString) {
		// Replacing all occurrences might be dangerous if not specified, 
		// but standard "replace" usually replaces all or requires a count.
		// The prompt says "string replace method", similar to Claude/Gemini.
		// Usually they replace one specific instance or all.
		// Let's assume we replace all for now, or just the first one?
		// My tool 'replace' replaces one by default.
		// Let's do one replacement to be safe and match the 'replace' tool behavior.
		newContent := strings.Replace(fileContent, req.OldString, req.NewString, 1)
		err = os.WriteFile(req.Path, []byte(newContent), 0644)
		if err != nil {
			return nil, err
		}
		return &replace_models.Response{Path: req.Path}, nil
	}

	// Fuzzy match
	bestMatchIdx := -1
	bestSimilarity := 0.0

	runes := []rune(fileContent)
	oldRunes := []rune(req.OldString)
	oldLen := len(oldRunes)

	if oldLen == 0 {
		return nil, api.NewError(api.BadRequest, "Old string cannot be empty")
	}
	
	for i := 0; i <= len(runes)-oldLen; i++ {
		windowRunes := runes[i : i+oldLen]
		distance := levenshtein(windowRunes, oldRunes)
		similarity := 1.0 - float64(distance)/float64(oldLen)
		
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatchIdx = i
		}
		
		if bestSimilarity >= 0.98 {
			break
		}
	}

	if bestSimilarity >= 0.98 {
		// Replace the best match
		prefix := string(runes[:bestMatchIdx])
		suffix := string(runes[bestMatchIdx+oldLen:])
		newContent := prefix + req.NewString + suffix
		
		err = os.WriteFile(req.Path, []byte(newContent), 0644)
		if err != nil {
			return nil, err
		}
		return &replace_models.Response{Path: req.Path}, nil
	}

	return nil, api.NewError(api.BadRequest, "Could not find a match with at least 98% similarity")
}

func levenshtein(r1, r2 []rune) int {
	len1 := len(r1)
	len2 := len(r2)

	column := make([]int, len1+1)
	for y := 1; y <= len1; y++ {
		column[y] = y
	}

	for x := 1; x <= len2; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= len1; y++ {
			oldkey := column[y]
			var incr int
			if r1[y-1] != r2[x-1] {
				incr = 1
			}
			column[y] = min(min(column[y]+1, column[y-1]+1), lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[len1]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
