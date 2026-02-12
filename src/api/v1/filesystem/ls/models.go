package ls

import (
	"agent-dev-environment/src/library/api"
	"time"
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

type FileInfo struct {
	Name        string    `json:"name"`
	IsDirectory bool      `json:"is_directory"`
	Size        int64     `json:"size"`
	Mode        string    `json:"mode,omitempty"`
	ModTime     time.Time `json:"mod_time,omitempty"`
}

type Response struct {
	Entries []FileInfo `json:"entries"`
}
