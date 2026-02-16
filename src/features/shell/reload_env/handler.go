package reload_env

import (
	"bufio"
	"os"
	"strings"

	"agent-dev-environment/src/api/v1"
	reload_models "agent-dev-environment/src/api/v1/shell/reload_env"
	"agent-dev-environment/src/library/api"
	"agent-dev-environment/src/library/logger"
	"context"
)

func Handler(ctx context.Context, req reload_models.Request) (*v1.EmptyResponse, error) {
	// Load the specified env file into the current process
	if err := loadDotEnv(req.Path); err != nil {
		logger.Error("Failed to load env file", "error", err, "path", req.Path)
		return nil, api.NewError(api.InternalServerError, "failed to load environment variables into process: "+err.Error())
	}

	return &v1.EmptyResponse{}, nil
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
