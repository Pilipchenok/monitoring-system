package service

import (
	"context"

	"monitoring-system/internal/model"
)

type MetricsStorage interface {
	SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error
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
