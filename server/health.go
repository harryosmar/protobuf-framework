package server

import (
	"encoding/json"
	"net/http"

	"github.com/harryosmar/protobuf-go/config"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	Status      string `json:"status"`
}

// HealthHandler returns the health status of the service
func HealthHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			ServiceName: cfg.AppName,
			Version:     cfg.AppVersion,
			Status:      "healthy",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
