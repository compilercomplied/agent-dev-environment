package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"agent-dev-environment/src/api/v1"
	"agent-dev-environment/src/library/logger"
)

// Status Codes as constants for clarity
const (
	OK                  = http.StatusOK
	Created             = http.StatusCreated
	BadRequest          = http.StatusBadRequest
	NotFound            = http.StatusNotFound
	Conflict            = http.StatusConflict
	InternalServerError = http.StatusInternalServerError
)

// AppError is a custom error type that carries an HTTP status code
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code int, message string) error {
	return &AppError{Code: code, Message: message}
}

// HandlerFunc is our "Clean Handler" signature
type HandlerFunc[Req any, Res any] func(req Req) (*Res, error)

// Validator interface for request structures
type Validator interface {
	Validate() error
}

// WrappedHandler converts a Clean Handler into a standard http.HandlerFunc
func WrappedHandler[Req any, Res any](hf HandlerFunc[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		// Only decode if there is a body. This allows GET or empty-body POSTs to work.
		if r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				respondError(w, "Invalid request body", BadRequest)
				return
			}
		}

		// Automatic validation if the request implements Validator
		if v, ok := any(req).(Validator); ok {
			if err := v.Validate(); err != nil {
				handleError(w, err)
				return
			}
		}

		res, err := hf(req)
		if err != nil {
			handleError(w, err)
			return
		}

		respond(w, res)
	}
}

func handleError(w http.ResponseWriter, err error) {
	var fErr *AppError
	if errors.As(err, &fErr) {
		respondError(w, fErr.Message, fErr.Code)
		return
	}

	// Fallback for unknown errors
	logger.Error("Unexpected error", "error", err)
	respondError(w, "Internal server error", InternalServerError)
}

func respond(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(OK)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v1.ErrorResponse{Error: message})
}