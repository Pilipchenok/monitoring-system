package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"monitoring-system/internal/model"
)

type ServiceStorage interface {
	SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error
	CleanOldMetrics(ctx context.Context, threshold time.Time) error
	SelectLastMetrics(ctx context.Context, count int) ([]model.Metric, error)
}

type Handler struct {
	service ServiceStorage
}

func NewHandler(service ServiceStorage) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SaveMetricsHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	countMetrics := 10
	metrics, err := h.service.SelectLastMetrics(r.Context(), countMetrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
