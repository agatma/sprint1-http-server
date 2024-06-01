package rest

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage"
	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/core/service"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
)

func TestHandler_SetMetricValueSuccess(t *testing.T) {
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
				Value: "13.12",
				Type:  domain.Gauge,
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
				Type:  domain.Counter,
			},
			method: http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, bytes.NewBuffer(make([]byte, 0)))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.metric.Name)
			rctx.URLParams.Add("metricType", tt.metric.Type)
			rctx.URLParams.Add("metricValue", tt.metric.Value)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			metricStorage, _ := storage.NewStorage(storage.Config{
				Memory: &memory.Config{},
			})
			metricService := service.NewMetricService(metricStorage)
			h := handler{
				metricService: metricService,
			}
			h.SetMetricValue(w, r)
			result := w.Result()
			defer func() {
				err := result.Body.Close()
				logger.Log.Error("error occurred during closing body", zap.Error(err))
			}()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			value, _ := h.metricService.GetMetricValue(tt.metric.Type, tt.metric.Name)
			assert.Equal(t, tt.metric.Value, value)
		})
	}
}

func TestHandler_SetMetricValueFailed(t *testing.T) {
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
				Type:  domain.Gauge,
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
				Type:  domain.Counter,
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
			r := httptest.NewRequest(http.MethodGet, tt.url, bytes.NewBuffer(make([]byte, 0)))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.metric.Name)
			rctx.URLParams.Add("metricType", tt.metric.Type)
			rctx.URLParams.Add("metricValue", tt.metric.Value)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			metricStorage, _ := storage.NewStorage(storage.Config{
				Memory: &memory.Config{},
			})
			metricService := service.NewMetricService(metricStorage)
			h := handler{
				metricService: metricService,
			}
			h.SetMetricValue(w, r)
			result := w.Result()
			defer func() {
				err := result.Body.Close()
				logger.Log.Error("error occurred during closing body", zap.Error(err))
			}()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
