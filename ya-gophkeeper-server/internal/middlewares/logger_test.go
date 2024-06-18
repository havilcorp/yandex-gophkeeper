// Package middleware мидлвар для логирования запросов
package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	})
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rw := httptest.NewRecorder()
	lm := LogMiddleware(testHandler)
	lm.ServeHTTP(rw, r)
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "OK" {
		t.Errorf("expected ABC got %v", string(data))
	}
}
