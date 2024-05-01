package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGaugeMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	type Metric struct {
		Name  string
		Value string
		Type  string
	}

	tests := []struct {
		name   string
		url    string
		metric Metric
		want   want
		method string
	}{
		{
			name: "statusOkGauge",
			url:  "/update/{metricType}/{metricName}/{metricValue}",
			metric: Metric{
				Name:  "someMetric",
				Value: "13.0",
				Type:  gauge,
			},
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "statusOkCounter",
			url:  "/update/{metricType}/{metricName}/{metricValue}",
			metric: Metric{
				Name:  "someMetric",
				Value: "13",
				Type:  counter,
			},
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "statusOkGauge",
			url:  "/update/{metricType}/{metricName}/{metricValue}",
			metric: Metric{
				Name:  "someMetric",
				Value: "13.0",
				Type:  "unknown",
			},
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name: "statusIncorrectMetricValue",
			url:  "/update/{metricType}/{metricName}/{metricValue}",
			metric: Metric{
				Name:  "someMetric",
				Value: "string",
				Type:  gauge,
			},
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name: "statusIncorrectMetricValue",
			url:  "/update/{metricType}/{metricName}/{metricValue}",
			metric: Metric{
				Name:  "someMetric",
				Value: "string",
				Type:  counter,
			},
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.metric.Name)
			rctx.URLParams.Add("metricType", tt.metric.Type)
			rctx.URLParams.Add("metricValue", tt.metric.Value)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			AddMetric(w, r)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
