package search

import (
	"agent-dev-environment/src/api/v1"
	"agent-dev-environment/src/api/v1/filesystem/search"
	"agent-dev-environment/src/library/api"
	"bytes"
	"os"
	"os/exec"
)

func Handler(req search.Request) (*v1.CommandResponse, error) {
	// First verify the path exists
	_, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "Path not found")
		}
		return nil, err
	}

	// Execute ripgrep (rg) command
	output, err := executeRipgrep(req)
	if err != nil {
		return nil, err
	}

	return &v1.CommandResponse{CommandOutput: output}, nil
}

func executeRipgrep(req search.Request) (string, error) {
	var args []string

	if req.FilesWithMatches {
		args = append(args, "--files-with-matches")
	}
	if req.IgnoreCase {
		args = append(args, "-i")
	}
	
	args = append(args, req.Pattern, req.Path)

	cmd := exec.Command("rg", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// rg returns exit code 1 if no matches are found, which is not an error for us
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "", nil
		}
		return "", api.NewError(api.InternalServerError, "rg command failed: "+stderr.String())
	}

	return stdout.String(), nil
}
