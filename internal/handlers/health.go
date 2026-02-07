package handlers

import (
	"net/http"

	results "github.com/GuyOz5252/go-app/pkg/api"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	results.ApiResponse(w, r, http.StatusOK, "healthy")
}
