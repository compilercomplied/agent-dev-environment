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

var allowedCommands = map[string]func([]string) error{
	"cargo":   func(args []string) error { return nil },
	"cat":     func(args []string) error { return nil },
	"cp":      func(args []string) error { return nil },
	"curl":    validateCurl,
	"diff":    func(args []string) error { return nil },
	"env":     func(args []string) error { return nil },
	"find":    func(args []string) error { return nil },
	"git":     func(args []string) error { return nil },
	"go":      func(args []string) error { return nil },
	"grep":    func(args []string) error { return nil },
	"head":    func(args []string) error { return nil },
	"jq":      func(args []string) error { return nil },
	"ls":      func(args []string) error { return nil },
	"make":    func(args []string) error { return nil },
	"mise":    func(args []string) error { return nil },
	"mkdir":   func(args []string) error { return nil },
	"mv":      func(args []string) error { return nil },
	"node":    func(args []string) error { return nil },
	"npm":     func(args []string) error { return nil },
	"npx":     func(args []string) error { return nil },
	"pip":     func(args []string) error { return nil },
	"ps":      func(args []string) error { return nil },
	"python":  func(args []string) error { return nil },
	"python3": func(args []string) error { return nil },
	"rg":      func(args []string) error { return nil },
	"rm":      func(args []string) error { return nil },
	"sed":     func(args []string) error { return nil },
	"echo":    func(args []string) error { return nil },
	"tee":     func(args []string) error { return nil },
	"base64":  func(args []string) error { return nil },
	"tail":    func(args []string) error { return nil },
	"touch":   func(args []string) error { return nil },
	"tsc":     func(args []string) error { return nil },
	"which":   func(args []string) error { return nil },
	"yarn":    func(args []string) error { return nil },
}

func (r Request) Validate() error {
	validator, ok := allowedCommands[r.Command]
	if !ok {
		return api.NewError(api.BadRequest, fmt.Sprintf("command '%s' is not allowed", r.Command))
	}

	if err := validator(r.Args); err != nil {
		return api.NewError(api.BadRequest, err.Error())
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
