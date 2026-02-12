package move

import (
	"agent-dev-environment/src/library/api"
)

type Request struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func (r Request) Validate() error {
	if r.Source == "" {
		return api.NewError(api.BadRequest, "Source path is required")
	}
	if r.Destination == "" {
		return api.NewError(api.BadRequest, "Destination path is required")
	}
	return nil
}
