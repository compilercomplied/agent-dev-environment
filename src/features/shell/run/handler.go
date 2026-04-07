package run

import (
	"agent-dev-environment/src/api/v1"
	"agent-dev-environment/src/api/v1/shell/run"
	"agent-dev-environment/src/library/api"
	"bytes"
	"context"
	"os/exec"
)

// Handler runs a shell command
// @Summary Run command
// @Description Executes a shell command and returns the output
// @Accept  json
// @Produce  json
// @Param   request  body  run.Request  true  "Run Command Request"
// @Success 200 {object} v1.CommandResponse
// @Router /api/v1/shell/run [post]
func Handler(ctx context.Context, req run.Request) (*v1.CommandResponse, error) {
	cmd := exec.CommandContext(ctx, req.Command, req.Args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, api.NewError(api.InternalServerError, "command failed: "+stderr.String())
	}

	return &v1.CommandResponse{CommandOutput: stdout.String()}, nil
}
