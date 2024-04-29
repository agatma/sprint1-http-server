package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
		method  string
	}{
		{
			name:    "statusOkCounter",
			request: "/update/counter/someMetric/527",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:    "statusNotFoundEmptyMetricType",
			request: "/update/counter/",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "statusIncorrectMetricValue",
			request: "/update/counter/someMetric/string",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "statusIncorrectMethod",
			request: "/update/counter/someMetric/string",
			method:  http.MethodDelete,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			CounterMetricHandler(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestGaugeMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
		method  string
	}{
		{
			name:    "statusOkCounter",
			request: "/update/gauge/someMetric/13.0",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:    "statusNotFoundEmptyMetricType",
			request: "/update/gauge/",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "statusIncorrectMetricValue",
			request: "/update/gauge/someMetric/string",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "statusIncorrectMethod",
			request: "/update/gauge/someMetric/string",
			method:  http.MethodDelete,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			GaugeMetricHandler(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
