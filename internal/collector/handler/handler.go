package handler

import (
	"encoding/json"
	"net/http"
	"context"

	"monitoring-system/internal/model"
)

type ServiceStorage interface {
	SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error
}

type Handler struct {
	service ServiceStorage
}

func NewHandler(service ServiceStorage) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AgentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data model.ServerMetrics
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.SaveMetrics(r.Context(), data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
