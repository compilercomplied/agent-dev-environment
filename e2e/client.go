package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"agent-dev-environment/src/api/v1"
	create_models "agent-dev-environment/src/api/v1/filesystem/create_file"
	delete_models "agent-dev-environment/src/api/v1/filesystem/delete"
	ls_models "agent-dev-environment/src/api/v1/filesystem/ls"
	mkdir_models "agent-dev-environment/src/api/v1/filesystem/mkdir"
	move_models "agent-dev-environment/src/api/v1/filesystem/move"
	read_models "agent-dev-environment/src/api/v1/filesystem/read"
	replace_models "agent-dev-environment/src/api/v1/filesystem/replace"
	search_models "agent-dev-environment/src/api/v1/filesystem/search"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (%d): %s", e.Status, e.Message)
}

func NewClient() *Client {
	baseURL := os.Getenv("E2E_SERVER_URL")
	if baseURL == "" {
		panic("E2E_SERVER_URL environment variable is required for e2e tests")
	}
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{},
	}
}

func (c *Client) CreateFile(req create_models.Request) (*v1.EmptyResponse, error) {
	return call[create_models.Request, v1.EmptyResponse](c, "POST", "/api/v1/filesystem/create_file", req)
}

func (c *Client) Mkdir(req mkdir_models.Request) (*v1.EmptyResponse, error) {
	return call[mkdir_models.Request, v1.EmptyResponse](c, "POST", "/api/v1/filesystem/mkdir", req)
}

func (c *Client) ReadFile(req read_models.Request) (*read_models.Response, error) {
	return call[read_models.Request, read_models.Response](c, "POST", "/api/v1/filesystem/read", req)
}

func (c *Client) DeleteFile(req delete_models.Request) (*v1.EmptyResponse, error) {
	return call[delete_models.Request, v1.EmptyResponse](c, "POST", "/api/v1/filesystem/delete", req)
}

func (c *Client) MoveFile(req move_models.Request) (*v1.EmptyResponse, error) {
	return call[move_models.Request, v1.EmptyResponse](c, "POST", "/api/v1/filesystem/move", req)
}

func (c *Client) ListFiles(req ls_models.Request) (*v1.CommandResponse, error) {
	return call[ls_models.Request, v1.CommandResponse](c, "POST", "/api/v1/filesystem/ls", req)
}

func (c *Client) Search(req search_models.Request) (*v1.CommandResponse, error) {
	return call[search_models.Request, v1.CommandResponse](c, "POST", "/api/v1/filesystem/search", req)
}

func (c *Client) Replace(req replace_models.Request) (*replace_models.Response, error) {
	return call[replace_models.Request, replace_models.Response](c, "POST", "/api/v1/filesystem/replace", req)
}

func call[Req any, Res any](c *Client, method, path string, payload Req) (*Res, error) {
	url := c.BaseURL + path
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errRes v1.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errRes); err != nil {
			return nil, &APIError{
				Status:  resp.StatusCode,
				Message: fmt.Sprintf("could not decode error response: %v", err),
			}
		}
		return nil, &APIError{
			Status:  resp.StatusCode,
			Message: errRes.Error,
		}
	}

	var res Res
	if resp.ContentLength != 0 {
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}
	}

	return &res, nil
}
