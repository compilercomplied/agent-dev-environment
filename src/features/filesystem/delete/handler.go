package delete

import (
	"errors"
	"os"

	"agent-dev-environment/src/api/v1"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	"agent-dev-environment/src/library/api"
)

func Handler(req delete_models.Request) (*v1.EmptyResponse, error) {
	stat, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "File or directory not found")
		}
		return nil, err
	}

	if stat.IsDir() && !req.Recursive {
		return nil, api.NewError(api.BadRequest, "Cannot delete directory without recursive flag")
	}

	if req.Recursive {
		err = os.RemoveAll(req.Path)
	} else {
		err = os.Remove(req.Path)
	}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, api.NewError(api.NotFound, "File or directory not found")
		}
		return nil, err
	}

	return &v1.EmptyResponse{}, nil
}
