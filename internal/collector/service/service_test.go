package service

import (
	"context"
	"errors"
	"testing"

	"monitoring-system/internal/model"
)

type MockStorage struct {
	mockErr error 
}

func (m *MockStorage) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {
	return m.mockErr
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
			mock := &MockStorage{mockErr: tt.mockErr}
			svc := NewService(mock)

			err := svc.SaveMetrics(context.Background(), model.ServerMetrics{})

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}