package create_file

import "agent-dev-environment/src/library/api"

type Request struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (r Request) Validate() error {
	if r.Path == "" {
		return api.NewError(api.BadRequest, "Path is required")
	}
	// Content can be empty, so no validation needed for it for now.
	return nil
}