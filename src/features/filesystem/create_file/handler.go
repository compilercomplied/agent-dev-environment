package create_file

import (
	"errors"
	"os"
	"path/filepath"

	"agent-dev-environment/src/api/v1"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	"agent-dev-environment/src/library/api"
	"context"
)

// Handler creates a new file
// @Summary Create file
// @Description Creates a new file with the specified content
// @Accept  json
// @Produce  json
// @Param   request  body  create_models.Request  true  "Create File Request"
// @Success 200 {object} v1.EmptyResponse
// @Router /api/v1/filesystem/create_file [post]
func Handler(ctx context.Context, req create_models.Request) (*v1.EmptyResponse, error) {
	dir := filepath.Dir(req.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	flags := os.O_WRONLY | os.O_CREATE
	if req.Overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(req.Path, flags, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, api.NewError(api.Conflict, "File already exists")
		}
		return nil, err
	}
	defer file.Close()

	if _, err := file.WriteString(req.Content); err != nil {
		return nil, err
	}

	return &v1.EmptyResponse{}, nil
}
