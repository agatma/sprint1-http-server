package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/agatma/sprint1-http-server/internal/compress"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, fmt.Errorf("failed to write response %w", err)
	}
	r.responseData.size += size
	return size, nil
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingRequestMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		logger.Log.Info("got incoming http request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", respData.status),
			zap.Int("size", respData.size),
			zap.String("duration", duration.String()),
		)
	}
	return http.HandlerFunc(logFn)
}

func CompressRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		cr, err := compress.NewCompressReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := cr.Close(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		r.Body = cr
		next.ServeHTTP(w, r)
	})
}

func CompressResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), `gzip`) {
			next.ServeHTTP(w, r)
			return
		}
		cw := compress.NewCompressWriter(w)
		defer func() {
			if err := cw.Close(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		w.Header().Set("Content-Encoding", `gzip`)

		next.ServeHTTP(cw, r)
	})
}
