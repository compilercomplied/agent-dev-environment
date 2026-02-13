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
	runes := []rune(fileContent)
	oldRunes := []rune(req.OldString)
	oldLen := len(oldRunes)

	if oldLen == 0 {
		return nil, api.NewError(api.BadRequest, "Old string cannot be empty")
	}

	// Find all matches with similarity >= 0.98
	type match struct {
		index      int
		similarity float64
	}
	var matches []match

	// Optimization: If exact matches exist, check their uniqueness first
	exactCount := strings.Count(fileContent, req.OldString)
	if exactCount > 1 {
		return nil, api.NewError(api.BadRequest, "Ambiguous replacement: multiple exact matches found. Please provide more context.")
	}

	// Slidding window for fuzzy matching
	for i := 0; i <= len(runes)-oldLen; i++ {
		windowRunes := runes[i : i+oldLen]
		distance := levenshtein(windowRunes, oldRunes)
		similarity := 1.0 - float64(distance)/float64(oldLen)

		if similarity >= 0.98 {
			matches = append(matches, match{index: i, similarity: similarity})
			// If we find more than one match, we can stop early if they are distinct enough.
			// However, slidding window might find the same match at adjacent indices.
			// Let's filter for distinct matches later.
		}
	}

	// Filter for distinct matches (matches that don't overlap significantly)
	var distinctMatches []match
	if len(matches) > 0 {
		distinctMatches = append(distinctMatches, matches[0])
		for i := 1; i < len(matches); i++ {
			lastMatch := distinctMatches[len(distinctMatches)-1]
			// If the new match starts after the last match ends, it's distinct
			if matches[i].index >= lastMatch.index+oldLen {
				distinctMatches = append(distinctMatches, matches[i])
			} else if matches[i].similarity > lastMatch.similarity {
				// If they overlap, keep the one with higher similarity
				distinctMatches[len(distinctMatches)-1] = matches[i]
			}
		}
	}

	if len(distinctMatches) == 0 {
		return nil, api.NewError(api.BadRequest, "Could not find a match with at least 98% similarity")
	}

	if len(distinctMatches) > 1 {
		return nil, api.NewError(api.BadRequest, "Ambiguous replacement: multiple matches found. Please provide more context to uniquely identify the target.")
	}

	// Perform the replacement
	bestMatch := distinctMatches[0]
	prefix := string(runes[:bestMatch.index])
	suffix := string(runes[bestMatch.index+oldLen:])
	newContent := prefix + req.NewString + suffix

	err = os.WriteFile(req.Path, []byte(newContent), 0644)
	if err != nil {
		return nil, err
	}
	return &replace_models.Response{Path: req.Path}, nil
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
