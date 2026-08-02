package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"monitoring-system/internal/model"
)

type MockStorage struct {
	mockSaveErr  error
	mockCleanErr error
}

func (m *MockStorage) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {
	return m.mockSaveErr
}

func (m *MockStorage) CleanOldMetrics(ctx context.Context, threshold time.Time) error {
	return m.mockCleanErr
}

func TestService_SaveMetrics(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		wantErr bool
	}{
		{
			name:    "Успешное сохранение метрик",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "Ошибка на уровне базы данных",
			mockErr: errors.New("db connection timeout"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockStorage{mockSaveErr: tt.mockErr}
			svc := NewService(mock)

			err := svc.SaveMetrics(context.Background(), model.ServerMetrics{})

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CleanOldMetrics(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		wantErr bool
	}{
		{
			name:    "Успешная очистка устаревших метрик",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "Ошибка БД при очистке",
			mockErr: errors.New("database locked"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockStorage{mockCleanErr: tt.mockErr}
			svc := NewService(mock)

			err := svc.CleanOldMetrics(context.Background(), time.Now())

			if (err != nil) != tt.wantErr {
				t.Errorf("CleanOldMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}