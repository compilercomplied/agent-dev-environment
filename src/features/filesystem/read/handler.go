package read

import (
	"os"

	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	"agent-dev-environment/src/library/api"
)

func Handler(req read_models.Request) (*read_models.Response, error) {
	content, err := os.ReadFile(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "File not found")
		}
		return nil, err
	}

	return &read_models.Response{Content: string(content)}, nil
}
