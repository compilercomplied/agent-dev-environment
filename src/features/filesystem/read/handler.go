package read

import (
	"bufio"
	"os"
	"strings"

	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"agent-dev-environment/src/library/api"
)

func Handler(req read_models.Request) (*read_models.Response, error) {
	file, err := os.Open(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "File not found")
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	
	currentLine := 0
	offset := 0
	if req.Offset != nil {
		offset = *req.Offset
	}
	
	limit := -1
	if req.Limit != nil {
		limit = *req.Limit
	}

	totalLines := 0
	for scanner.Scan() {
		if currentLine >= offset && (limit == -1 || len(lines) < limit) {
			lines = append(lines, scanner.Text())
		}
		currentLine++
		totalLines++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if offset >= totalLines && totalLines > 0 {
		return nil, api.NewError(api.BadRequest, "Offset is out of bounds")
	}

	hasMore := false
	if limit != -1 && offset+limit < totalLines {
		hasMore = true
	}

	return &read_models.Response{
		Content:    strings.Join(lines, "\n"),
		TotalLines: totalLines,
		HasMore:    hasMore,
		LinesRead:  len(lines),
	}, nil
}
