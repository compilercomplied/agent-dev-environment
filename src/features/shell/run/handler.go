package run

import (
	"agent-dev-environment/src/api/v1"
	"agent-dev-environment/src/api/v1/shell/run"
	"agent-dev-environment/src/library/api"
	"bytes"
	"context"
	"os/exec"
)

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
