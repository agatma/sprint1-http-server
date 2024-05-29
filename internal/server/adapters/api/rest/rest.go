package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/api/middleware"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
)

const (
	metricType  = "metricType"
	metricValue = "metricValue"
	metricName  = "metricName"
)

type MetricService interface {
	GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse
	SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse
	GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse
}

type handler struct {
	metricService MetricService
}

type API struct {
	srv *http.Server
}

func (a *API) Run() error {
	if err := a.srv.ListenAndServe(); err != nil {
		logger.Log.Error("error occured during running server: ", zap.Error(err))
		return fmt.Errorf("failed run server: %w", err)
	}
	return nil
}

func NewAPI(metricService MetricService, cfg *Config) *API {
	h := &handler{
		metricService: metricService,
	}
	r := chi.NewRouter()
	r.Use(middleware.LoggingRequestMiddleware)
	r.Use(middleware.CompressRequestMiddleware)
	r.Use(middleware.CompressResponseMiddleware)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.SetMetricValueJSON)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.SetMetricValue)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricValueJSON)
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

func (h *handler) SetMetricValue(w http.ResponseWriter, req *http.Request) {
	mType := chi.URLParam(req, metricType)
	mName := chi.URLParam(req, metricName)
	mValue := chi.URLParam(req, metricValue)
	response := h.metricService.SetMetricValue(&domain.SetMetricRequest{
		MetricType:  mType,
		MetricName:  mName,
		MetricValue: mValue,
	})
	if response.Error != nil {
		logger.Log.Error("failed to set metric",
			zap.String(metricValue, mValue),
			zap.String(metricType, mType),
			zap.String(metricName, mName),
			zap.Error(response.Error),
		)
		switch {
		case errors.Is(response.Error, domain.ErrIncorrectMetricType):
			http.Error(w, response.Error.Error(), http.StatusBadRequest)
		case errors.Is(response.Error, domain.ErrIncorrectMetricValue):
			http.Error(w, response.Error.Error(), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) SetMetricValueJSON(w http.ResponseWriter, req *http.Request) {
	var (
		request domain.Metrics
		mValue  string
	)
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
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
	switch request.MType {
	case domain.Gauge:
		mValue = strconv.FormatFloat(*request.Value, 'f', -1, 64)
	case domain.Counter:
		mValue = strconv.Itoa(int(*request.Delta))
	}
	serviceResponse := h.metricService.SetMetricValue(&domain.SetMetricRequest{
		MetricType:  request.MType,
		MetricName:  request.ID,
		MetricValue: mValue,
	})

	if serviceResponse.Error != nil {
		logger.Log.Error("failed to set metric",
			zap.String(metricValue, mValue),
			zap.String(metricType, request.MType),
			zap.String(metricName, request.ID),
			zap.Error(serviceResponse.Error),
		)
		switch {
		case errors.Is(serviceResponse.Error, domain.ErrIncorrectMetricType):
			http.Error(w, serviceResponse.Error.Error(), http.StatusBadRequest)
		case errors.Is(serviceResponse.Error, domain.ErrIncorrectMetricValue):
			http.Error(w, serviceResponse.Error.Error(), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	metricResponse := h.metricService.GetMetricValue(&domain.MetricRequest{
		MetricType: request.MType,
		MetricName: request.ID,
	})
	response, err := createMetricResponse(&request, metricResponse.MetricValue)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}

func (h *handler) GetMetricValue(w http.ResponseWriter, req *http.Request) {
	mType, mName := chi.URLParam(req, metricType), chi.URLParam(req, metricName)
	response := h.metricService.GetMetricValue(&domain.MetricRequest{
		MetricType: mType,
		MetricName: mName,
	})
	if !response.Found {
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	if response.Error != nil {
		logger.Log.Error("failed to get metric",
			zap.String(metricType, mType),
			zap.String(metricName, mName),
			zap.Error(response.Error),
		)
		if errors.Is(response.Error, domain.ErrIncorrectMetricType) {
			http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if _, err := w.Write([]byte(response.MetricValue)); err != nil {
		return
	}
}

func (h *handler) GetMetricValueJSON(w http.ResponseWriter, req *http.Request) {
	var request domain.Metrics
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricResponse := h.metricService.GetMetricValue(&domain.MetricRequest{
		MetricType: request.MType,
		MetricName: request.ID,
	})
	if !metricResponse.Found {
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	if metricResponse.Error != nil {
		logger.Log.Error("failed to get metric",
			zap.String(metricType, request.MType),
			zap.String(metricName, request.ID),
			zap.Error(metricResponse.Error),
		)
		if errors.Is(metricResponse.Error, domain.ErrIncorrectMetricType) {
			http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	response, err := createMetricResponse(&request, metricResponse.MetricValue)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
}

func (h *handler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	gauge := h.metricService.GetAllMetrics(&domain.GetAllMetricsRequest{MetricType: domain.Gauge})
	if gauge.Error != nil {
		logger.Log.Error(
			"failed to get an item",
			zap.String(metricType, domain.Gauge),
			zap.Error(gauge.Error),
		)
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	counter := h.metricService.GetAllMetrics(&domain.GetAllMetricsRequest{MetricType: domain.Counter})
	if counter.Error != nil {
		logger.Log.Error(
			"failed to get an item",
			zap.String(metricType, domain.Counter),
			zap.Error(counter.Error),
		)
		http.Error(w, domain.ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}
	html := "<html><body><ul>"
	for key, value := range gauge.Values {
		html += fmt.Sprintf("<li>gauge: %s: %v</li>", key, value)
	}
	for key, value := range counter.Values {
		html += fmt.Sprintf("<li>counter: %s: %v</li>", key, value)
	}
	html += "</ul></body></html>"
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func createMetricResponse(request *domain.Metrics, value string) (*domain.Metrics, error) {
	var response domain.Metrics
	switch request.MType {
	case domain.Gauge:
		gaugeValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value: %s, error: %w", value, err)
		}
		response = domain.Metrics{
			ID:    request.ID,
			MType: request.MType,
			Value: &gaugeValue,
		}
	case domain.Counter:
		counterValue, err := strconv.Atoi(value)
		counterInt64Value := int64(counterValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value: %s, error: %w", value, err)
		}
		response = domain.Metrics{
			ID:    request.ID,
			MType: request.MType,
			Delta: &counterInt64Value,
		}
	}
	return &response, nil
}
