// Package middleware мидлвар для логирования запросов
package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	responseData struct {
		status int
		size   int
	}
)

type (
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LogMiddleware мидлвар для логирования запросов
func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 200,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		logrus.Infof("%s %s (%d) %s %d byte", r.Method, r.RequestURI, responseData.status, duration, responseData.size)
	})
}
