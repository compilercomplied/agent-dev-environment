package move

import (
	"agent-dev-environment/src/library/api"
	v1 "agent-dev-environment/src/api/v1"
	models "agent-dev-environment/src/api/v1/filesystem/move"
	"os"
)

func Handler(req models.Request) (*v1.EmptyResponse, error) {
	// Check if source exists
	if _, err := os.Stat(req.Source); os.IsNotExist(err) {
		return nil, api.NewError(api.NotFound, "Source path does not exist")
	}

	// Check if destination already exists
	if _, err := os.Stat(req.Destination); err == nil {
		return nil, api.NewError(api.Conflict, "Destination path already exists")
	}

	// Perform the move operation using os.Rename
	if err := os.Rename(req.Source, req.Destination); err != nil {
		return nil, api.NewError(api.InternalServerError, "Failed to move file: "+err.Error())
	}

	return &v1.EmptyResponse{}, nil
}
