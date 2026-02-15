package chdir

import (
	"os"
	"path/filepath"

	"agent-dev-environment/src/api/v1"
	chdir_models "agent-dev-environment/src/api/v1/filesystem/chdir"
	"agent-dev-environment/src/library/api"
)

func Handler(req chdir_models.Request) (*v1.EmptyResponse, error) {
	absPath, err := filepath.Abs(req.Path)
	if err != nil {
		return nil, api.NewError(api.BadRequest, "Invalid path: "+err.Error())
	}

	if err := os.Chdir(absPath); err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "Directory not found")
		}
		return nil, api.NewError(api.InternalServerError, "Failed to change directory: "+err.Error())
	}

	return &v1.EmptyResponse{}, nil
}
