package replace

import "agent-dev-environment/src/library/api"

type Request struct {
	Path                 string `json:"path"`
	OldString            string `json:"old_string"`
	NewString            string `json:"new_string"`
}

func (r Request) Validate() error {
	if r.Path == "" {
		return api.NewError(api.BadRequest, "Path is required")
	}
	if r.OldString == "" {
		return api.NewError(api.BadRequest, "Old string is required")
	}
	// NewString can be empty (for deletion)
	return nil
}

type Response struct {
	Path string `json:"path"`
}
