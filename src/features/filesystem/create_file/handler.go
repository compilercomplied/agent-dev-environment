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
