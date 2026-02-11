package main

import (
	"agent-dev-environment/src/internal/middleware"
	"agent-dev-environment/src/features/filesystem/create_file"
	"agent-dev-environment/src/features/filesystem/read"
	"agent-dev-environment/src/library/api"
	"agent-dev-environment/src/library/config"
	"agent-dev-environment/src/library/logger"
	"net/http"
)

func main() {
	logFormat := config.GetValue("LOGGING_TYPE")
	logger.Init(logFormat)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/v1/filesystem/read", api.WrappedHandler(read.Handler))
	mux.HandleFunc("POST /api/v1/filesystem/create_file", api.WrappedHandler(create_file.Handler))

	handler := middleware.PanicRecovery(mux)

	port := "8080"
	logger.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
