package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"monitoring-system/internal/model"
)

type ServiceStorage interface {
	SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error
	CleanOldMetrics(ctx context.Context, threshold time.Time) error
	SelectLastMetrics(ctx context.Context, hostID int, count int) ([]model.Metric, error)
	GetAllHosts(ctx context.Context) ([]model.Host, error)
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
	hostIDStr := r.URL.Query().Get("host_id")
	if hostIDStr == "" {
		http.Error(w, "missing host_id parameter", http.StatusBadRequest)
		return
	}
	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		http.Error(w, "invalid host_id parameter", http.StatusBadRequest)
		return
	}

	countMetrics := 80
	metrics, err := h.service.SelectLastMetrics(r.Context(), hostID, countMetrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) GetHostsHandler(w http.ResponseWriter, r *http.Request) {
	hosts, err := h.service.GetAllHosts(r.Context())
	if err != nil {
		http.Error(w, "Failed to get hosts", http.StatusInternalServerError)
		return
	}

	if hosts == nil {
		hosts = []model.Host{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hosts)
}
