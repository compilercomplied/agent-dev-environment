package reload_env

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"agent-dev-environment/src/api/v1"
	"agent-dev-environment/src/library/api"
	"agent-dev-environment/src/library/logger"
)

func Handler(req v1.EmptyResponse) (*v1.CommandResponse, error) {
	// Execute mise run reload-env
	cmd := exec.Command("mise", "run", "reload-env")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, api.NewError(api.InternalServerError, "reload-env failed: "+stderr.String())
	}

	// After running the script, we load the .env file into the current process
	// so that subsequent shell commands inherit these variables.
	if err := loadDotEnv(".env"); err != nil {
		logger.Error("Failed to load .env file after reload-env", "error", err)
		// We don't return error here because the command itself succeeded
	}

	return &v1.CommandResponse{CommandOutput: stdout.String()}, nil
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]
		
		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
