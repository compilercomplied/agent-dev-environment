package search

import (
	"agent-dev-environment/src/library/api"
)

type Request struct {
	Path             string `json:"path"`
	Pattern          string `json:"pattern"`
	FilesWithMatches bool   `json:"files_with_matches"`
	IgnoreCase       bool   `json:"ignore_case"`
}

func (r Request) Validate() error {
	if r.Path == "" {
		return api.NewError(api.BadRequest, "Path is required")
	}
	if r.Pattern == "" {
		return api.NewError(api.BadRequest, "Pattern is required")
	}
	return nil
}
