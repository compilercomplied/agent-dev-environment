package getwd

import (
	"os"

	"agent-dev-environment/src/api/v1"
	getwd_models "agent-dev-environment/src/api/v1/filesystem/getwd"
	"agent-dev-environment/src/library/api"
)

func Handler(req v1.EmptyResponse) (*getwd_models.Response, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, api.NewError(api.InternalServerError, "Failed to get working directory: "+err.Error())
	}
	return &getwd_models.Response{Path: wd}, nil
}
