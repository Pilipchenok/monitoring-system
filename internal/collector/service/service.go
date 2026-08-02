package service

import (
	"context"
	"time"

	"monitoring-system/internal/model"
)

type MetricsStorage interface {
	SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error
	CleanOldMetrics(ctx context.Context, threshold time.Time) error
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
