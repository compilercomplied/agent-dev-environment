package main

import (
	"agent-dev-environment/src/internal/middleware"
	"agent-dev-environment/src/features/filesystem/create_file"
	"agent-dev-environment/src/features/filesystem/read"
	"agent-dev-environment/src/library/api"
	"agent-dev-environment/src/library/config"
	"agent-dev-environment/src/library/logger"
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	logFormat := config.GetValue("LOGGING_TYPE")
	logger.Init(logFormat)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", helloHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/v1/filesystem/read", api.WrappedHandler(read.Handler))
	mux.HandleFunc("POST /api/v1/filesystem/create_file", api.WrappedHandler(create_file.Handler))

	// Apply middleware
	handler := middleware.PanicRecovery(mux)

	port := "8080"
	logger.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
