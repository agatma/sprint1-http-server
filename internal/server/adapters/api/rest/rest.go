package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/api/middleware"
	"github.com/agatma/sprint1-http-server/internal/server/config"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
)

const (
	metricType  = "metricType"
	metricValue = "metricValue"
	metricName  = "metricName"
)

type MetricService interface {
	GetMetric(mType, mName string) (*domain.Metric, error)
	GetMetricValue(mType, mName string) (string, error)
	SetMetric(m *domain.Metric) (*domain.Metric, error)
	SetMetricValue(m *domain.SetMetricRequest) (*domain.Metric, error)
	GetAllMetrics() (domain.MetricsList, error)
}

type handler struct {
	metricService MetricService
}

type API struct {
	srv *http.Server
}

func (a *API) Run() error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		if err := a.srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("server shutdown: ", zap.Error(err))
		}
	}()
	if err := a.srv.ListenAndServe(); err != nil {
		logger.Log.Error("error occurred during running server: ", zap.Error(err))
		return fmt.Errorf("failed run server: %w", err)
	}
	return nil
}

func NewAPI(metricService MetricService, cfg *config.Config) *API {
	h := &handler{
		metricService: metricService,
	}
	r := chi.NewRouter()
	r.Use(middleware.LoggingRequestMiddleware)
	r.Use(middleware.CompressRequestMiddleware)
	r.Use(middleware.CompressResponseMiddleware)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.SetMetric)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.SetMetricValue)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetric)
		r.Get("/{metricType}/{metricName}", h.GetMetricValue)
	})
	r.Get("/", h.GetAllMetrics)
	return &API{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      r,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
}

func handleSetMetricError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrItemNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrIncorrectMetricType) || errors.Is(err, domain.ErrIncorrectMetricValue):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetMetricError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrIncorrectMetricType) || errors.Is(err, domain.ErrItemNotFound) {
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *handler) SetMetricValue(w http.ResponseWriter, req *http.Request) {
	mType := chi.URLParam(req, metricType)
	mName := chi.URLParam(req, metricName)
	mValue := chi.URLParam(req, metricValue)
	_, err := h.metricService.SetMetricValue(&domain.SetMetricRequest{
		ID:    mName,
		MType: mType,
		Value: mValue,
	})
	if err != nil {
		logger.Log.Error("failed to set metric",
			zap.String(metricValue, mValue),
			zap.String(metricType, mType),
			zap.String(metricName, mName),
			zap.Error(err),
		)
		handleSetMetricError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) SetMetric(w http.ResponseWriter, req *http.Request) {
	var m domain.Metric
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := io.Copy(io.Discard, req.Body)
	if err != nil {
		logger.Log.Info("cannot read body", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = req.Body.Close()
	if err != nil {
		logger.Log.Info("cannot close body", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	metric, err := h.metricService.SetMetric(&m)

	if err != nil {
		logger.Log.Error("failed to set metric", zap.Error(err))
		handleSetMetricError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}

func (h *handler) GetMetricValue(w http.ResponseWriter, req *http.Request) {
	mType, mName := chi.URLParam(req, metricType), chi.URLParam(req, metricName)
	metricValue, err := h.metricService.GetMetricValue(mType, mName)
	if err != nil {
		logger.Log.Error("failed to get metric",
			zap.String(metricType, mType),
			zap.String(metricName, mName),
			zap.Error(err),
		)
		handleGetMetricError(w, err)
		return
	}
	if _, err := w.Write([]byte(metricValue)); err != nil {
		return
	}
}

func (h *handler) GetMetric(w http.ResponseWriter, req *http.Request) {
	var m domain.Metric
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, err := h.metricService.GetMetric(m.MType, m.ID)

	if err != nil {
		logger.Log.Error("failed to get metric", zap.Error(err))
		handleGetMetricError(w, err)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}

func (h *handler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	metrics, err := h.metricService.GetAllMetrics()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logger.Log.Error("failed to get all metrics", zap.Error(err))
		return
	}
	html := "<html><body><ul>"
	for _, metric := range metrics {
		switch metric.MType {
		case domain.Gauge:
			html += fmt.Sprintf("<li>mType: %s, mName: %s, Value %v", metric.MType, metric.ID, *metric.Value)
		case domain.Counter:
			html += fmt.Sprintf("<li>mType: %s, mName: %s, Value %v", metric.MType, metric.ID, *metric.Delta)
		}
	}
	html += "</ul></body></html>"
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
}
