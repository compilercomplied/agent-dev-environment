package mkdir

import (
	"os"

	"agent-dev-environment/src/api/v1"
	mkdir_models "agent-dev-environment/src/api/v1/filesystem/mkdir"
	"agent-dev-environment/src/library/api"
	"context"
)

// Handler creates a new directory
// @Summary Create directory
// @Description Creates a new directory at the specified path
// @Accept  json
// @Produce  json
// @Param   request  body  mkdir_models.Request  true  "Create Directory Request"
// @Success 200 {object} v1.EmptyResponse
// @Router /api/v1/filesystem/mkdir [post]
func Handler(ctx context.Context, req mkdir_models.Request) (*v1.EmptyResponse, error) {
	stat, err := os.Stat(req.Path)
	if err == nil {
		if !stat.IsDir() {
			return nil, api.NewError(api.Conflict, "Path already exists and is not a directory")
		}
		return &v1.EmptyResponse{}, nil
	}

	if err := os.MkdirAll(req.Path, 0755); err != nil {
		return nil, err
	}
	return &v1.EmptyResponse{}, nil
}
