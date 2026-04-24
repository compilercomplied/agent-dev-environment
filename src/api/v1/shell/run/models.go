package run

import (
	"agent-dev-environment/src/library/api"
	"fmt"
	"strings"
)

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func (r Request) Validate() error {
	if r.Command == "" {
		return api.NewError(api.BadRequest, "command is required")
	}

	if r.Command == "curl" {
		if err := validateCurl(r.Args); err != nil {
			return api.NewError(api.BadRequest, err.Error())
		}
	}

	return nil
}

func validateCurl(args []string) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			isLocalhost := strings.Contains(arg, "localhost") ||
				strings.Contains(arg, "127.0.0.1") ||
				strings.Contains(arg, "[::1]")

			if !isLocalhost {
				return fmt.Errorf("curl is restricted to localhost targets due to security reasons")
			}
		}
	}
	return nil
}
