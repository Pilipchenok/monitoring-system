package service

import (
	"context"
	"time"

	"monitoring-system/internal/model"
)

type MetricsStorage interface {
	SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error
	CleanOldMetrics(ctx context.Context, threshold time.Time) error
	SelectLastMetrics(ctx context.Context, hostID int, count int) ([]model.Metric, error)
	GetAllHosts(ctx context.Context) ([]model.Host, error)
}

type Service struct {
	storage MetricsStorage
}

func NewService(storage MetricsStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {
	return s.storage.SaveMetrics(ctx, metrics)
}

func (s *Service) CleanOldMetrics(ctx context.Context, threshold time.Time) error {
	return s.storage.CleanOldMetrics(ctx, threshold)
}

func (s *Service) SelectLastMetrics(ctx context.Context, hostID int, count int) ([]model.Metric, error) {
	return s.storage.SelectLastMetrics(ctx, hostID, count)
}

func (s *Service) GetAllHosts(ctx context.Context) ([]model.Host, error) {
	return s.storage.GetAllHosts(ctx)
}
