package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"monitoring-system/internal/model"
)

type MockService struct {
	mockSaveErr error
	mockSelectErr error
	mockData []model.Metric
}

func (m *MockService) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {
	return m.mockSaveErr
}

func (m *MockService) CleanOldMetrics(ctx context.Context, threshold time.Time) error {
	return nil
}

func (m *MockService) SelectLastMetrics(ctx context.Context, count int) ([]model.Metric, error) {
	return m.mockData, m.mockSelectErr
}

func TestHandler_SaveMetricsHandler(t *testing.T) {
	validJSON := []byte(`{"hostname":"test-mac","metrics":[{"name":"CPU","value":50.0}]}`)
	invalidJSON := []byte(`{"hostname":"test-mac", wrong format}`)

	tests := []struct {
		name string
		method string
		body []byte
		mockErr error
		expectedCode int
	}{
		{
			name: "Успешное сохранение",
			method: http.MethodPost,
			body: validJSON,
			mockErr: nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Невалидный JSON",
			method: http.MethodPost,
			body: invalidJSON,
			mockErr: nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Ошибка базы данных",
			method: http.MethodPost,
			body: validJSON,
			mockErr: errors.New("db error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{mockSaveErr: tt.mockErr}
			h := NewHandler(mockSvc)

			req, _ := http.NewRequest(tt.method, "/api/metrics", bytes.NewBuffer(tt.body))
			rr := httptest.NewRecorder()

			h.SaveMetricsHandler(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("SaveMetricsHandler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}
		})
	}
}

func TestHandler_GetMetricsHandler(t *testing.T) {
	tests := []struct {
		name string
		method string
		mockData []model.Metric
		mockErr error
		expectedCode int
	}{
		{
			name: "Успешное получение",
			method: http.MethodGet,
			mockData: []model.Metric{
				{Name: "CPU", Value: 42.5},
				{Name: "RAM", Value: 70.1},
			},
			mockErr: nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Ошибка базы данных",
			method: http.MethodGet,
			mockData: nil,
			mockErr: errors.New("db error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{
				mockSelectErr: tt.mockErr,
				mockData: tt.mockData,
			}
			h := NewHandler(mockSvc)

			req, _ := http.NewRequest(tt.method, "/api/metrics/latest", nil)
			rr := httptest.NewRecorder()

			h.GetMetricsHandler(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("GetMetricsHandler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusOK {
				var responseData []model.Metric
				err := json.NewDecoder(rr.Body).Decode(&responseData)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(responseData, tt.mockData) {
					t.Errorf("GetMetricsHandler returned wrong body: got %v want %v",
						responseData, tt.mockData)
				}

				if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
					t.Errorf("GetMetricsHandler returned wrong Content-Type: got %v want %v",
						contentType, "application/json")
				}
			}
		})
	}
}
