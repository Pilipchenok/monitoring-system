package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"monitoring-system/internal/model"
)

type MockStorage struct {
	mockSaveErr error
	mockCleanErr error
	mockSelectErr error
	mockSelectData []model.Metric
}

func (m *MockStorage) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {
	return m.mockSaveErr
}

func (m *MockStorage) CleanOldMetrics(ctx context.Context, threshold time.Time) error {
	return m.mockCleanErr
}

func (m *MockStorage) SelectLastMetrics(ctx context.Context, count int) ([]model.Metric, error) {
	return m.mockSelectData, m.mockSelectErr
}

func TestService_SaveMetrics(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		wantErr bool
	}{
		{
			name: "Успешное сохранение метрик",
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

func TestService_SelectLastMetrics(t *testing.T) {
	tests := []struct {
		name string
		count int
		mockData []model.Metric
		mockErr error
		wantErr bool
	}{
		{
			name:  "Успешное получение метрик",
			count: 2,
			mockData: []model.Metric{
				{Name: "CPU", Value: 45.5},
				{Name: "RAM", Value: 60.2},
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:     "Ошибка БД при получении",
			count:    10,
			mockData: nil,
			mockErr:  errors.New("database connection lost"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockStorage{
				mockSelectErr:  tt.mockErr,
				mockSelectData: tt.mockData,
			}
			svc := NewService(mock)

			data, err := svc.SelectLastMetrics(context.Background(), tt.count)

			if (err != nil) != tt.wantErr {
				t.Errorf("SelectLastMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(data, tt.mockData) {
				t.Errorf("SelectLastMetrics() got = %v, want %v", data, tt.mockData)
			}
		})
	}
}
