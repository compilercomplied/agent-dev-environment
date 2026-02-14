package read

import "agent-dev-environment/src/library/api"

type Request struct {
	Path   string `json:"path"`
	Offset *int   `json:"offset,omitempty"` // 0-based starting line
	Limit  *int   `json:"limit,omitempty"`  // Number of lines to read
}

func (r Request) Validate() error {
	if r.Path == "" {
		return api.NewError(api.BadRequest, "Path is required")
	}
	if r.Offset != nil && *r.Offset < 0 {
		return api.NewError(api.BadRequest, "Offset cannot be negative")
	}
	if r.Limit != nil && *r.Limit <= 0 {
		return api.NewError(api.BadRequest, "Limit must be greater than 0")
	}
	return nil
}

type Response struct {
	Content     string `json:"content"`
	TotalLines  int    `json:"total_lines"`
	HasMore     bool   `json:"has_more"`
	LinesRead   int    `json:"lines_read"`
}