package ls

import (
	"agent-dev-environment/src/library/api"
)

type Request struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
	Long      bool   `json:"long"`
}

func (r Request) Validate() error {
	if r.Path == "" {
		return api.NewError(api.BadRequest, "Path is required")
	}
	return nil
}

type Response struct {
	CommandOutput string `json:"command output"`
}
