package ls

import (
	"agent-dev-environment/src/api/v1/filesystem/ls"
	"agent-dev-environment/src/library/api"
	"bytes"
	"os"
	"os/exec"
)

func Handler(req ls.Request) (*ls.Response, error) {
	// First verify the path exists
	_, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "Path not found")
		}
		return nil, err
	}

	// Execute Linux ls command
	output, err := executeLinuxLS(req.Path, req.Recursive, req.Long)
	if err != nil {
		return nil, err
	}

	return &ls.Response{CommandOutput: output}, nil
}

func executeLinuxLS(path string, recursive bool, long bool) (string, error) {
	var args []string

	// Build ls command arguments
	if long {
		args = append(args, "-l")
	}
	if recursive {
		args = append(args, "-R")
	}
	args = append(args, path)

	cmd := exec.Command("ls", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", api.NewError(api.InternalServerError, "ls command failed: "+stderr.String())
	}

	return stdout.String(), nil
}

