package ls

import (
	"agent-dev-environment/src/api/v1"
	"agent-dev-environment/src/api/v1/filesystem/ls"
	"agent-dev-environment/src/library/api"
	"bytes"
	"context"
	"os"
	"os/exec"
)

// Handler lists directory contents
// @Summary List directory
// @Description Returns a list of files and directories in the specified path
// @Accept  json
// @Produce  json
// @Param   request  body  ls.Request  true  "List Directory Request"
// @Success 200 {object} v1.CommandResponse
// @Router /api/v1/filesystem/ls [post]
func Handler(ctx context.Context, req ls.Request) (*v1.CommandResponse, error) {
	// First verify the path exists
	_, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "Path not found")
		}
		return nil, err
	}

	// Execute Linux ls command
	output, err := executeLinuxLS(ctx, req.Path, req.Recursive, req.Long, req.Hidden)
	if err != nil {
		return nil, err
	}

	return &v1.CommandResponse{CommandOutput: output}, nil
}

func executeLinuxLS(ctx context.Context, path string, recursive bool, long bool, hidden bool) (string, error) {
	var args []string

	// Build ls command arguments
	if hidden {
		args = append(args, "-a")
	}
	if long {
		args = append(args, "-l")
	}
	if recursive {
		args = append(args, "-R")
	}
	args = append(args, path)

	cmd := exec.CommandContext(ctx, "ls", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", api.NewError(api.InternalServerError, "ls command failed: "+stderr.String())
	}

	return stdout.String(), nil
}

