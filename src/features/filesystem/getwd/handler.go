package getwd

import (
	"os"

	"agent-dev-environment/src/api/v1"
	getwd_models "agent-dev-environment/src/api/v1/filesystem/getwd"
	"agent-dev-environment/src/library/api"
	"context"
)

// Handler gets current working directory
// @Summary Get working directory
// @Description Returns the current working directory of the server
// @Accept  json
// @Produce  json
// @Param   request  body  v1.EmptyResponse  true  "Get Working Directory Request"
// @Success 200 {object} getwd_models.Response
// @Router /api/v1/filesystem/getwd [post]
func Handler(ctx context.Context, req v1.EmptyResponse) (*getwd_models.Response, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, api.NewError(api.InternalServerError, "Failed to get working directory: "+err.Error())
	}
	return &getwd_models.Response{Path: wd}, nil
}
