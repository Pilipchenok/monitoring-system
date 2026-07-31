package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"monitoring-system/internal/model"
)

type MockService struct {
	mockErr error
}

func (m *MockService) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {
	return m.mockErr
}

func TestHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		mockServiceErr error
		expectedStatus int
	}{
		{
			name:           "Успешный запрос (200 OK)",
			method:         http.MethodPost,
			body:           `{"hostname":"backend1","metrics":[{"name":"cpu","value":50}]}`,
			mockServiceErr: nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Неверный метод (405 Method Not Allowed)",
			method:         http.MethodGet,
			body:           ``,
			mockServiceErr: nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Кривой JSON (400 Bad Request)",
			method:         http.MethodPost,
			body:           `{"hostname":"backend1", broken json}`,
			mockServiceErr: nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Ошибка внутри сервиса/БД (500 Internal Error)",
			method:         http.MethodPost,
			body:           `{"hostname":"backend1","metrics":[{"name":"cpu","value":50}]}`,
			mockServiceErr: errors.New("database is dead"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{mockErr: tt.mockServiceErr}
			h := NewHandler(mockSvc)

			req := httptest.NewRequest(tt.method, "/", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}