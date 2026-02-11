package read

import "agent-dev-environment/src/library/api"

type Request struct {
	Path string `json:"path"`
}

func (r Request) Validate() error {
	if r.Path == "" {
		return api.NewError(api.BadRequest, "Path is required")
	}
	return nil
}

type Response struct {
	Content string `json:"content"`
}