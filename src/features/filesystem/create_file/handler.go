package create_file

import (
	"errors"
	"os"
	"path/filepath"

	"agent-dev-environment/src/api/v1"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	"agent-dev-environment/src/library/api"
)

func Handler(req create_models.Request) (*v1.EmptyResponse, error) {
	dir := filepath.Dir(req.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(req.Path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
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
