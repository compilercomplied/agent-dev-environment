package move

import (
	"agent-dev-environment/src/library/api"
	"context"
	v1 "agent-dev-environment/src/api/v1"
	models "agent-dev-environment/src/api/v1/filesystem/move"
	"os"
)

// Handler moves a file or directory
// @Summary Move file/directory
// @Description Moves the specified source to the destination path
// @Accept  json
// @Produce  json
// @Param   request  body  models.Request  true  "Move Request"
// @Success 200 {object} v1.EmptyResponse
// @Router /api/v1/filesystem/move [post]
func Handler(ctx context.Context, req models.Request) (*v1.EmptyResponse, error) {
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
