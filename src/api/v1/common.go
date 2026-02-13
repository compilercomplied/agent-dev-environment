package v1

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmptyResponse struct{}

type CommandResponse struct {
	CommandOutput string `json:"command_output"`
}