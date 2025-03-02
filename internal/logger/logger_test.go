package logger

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	assert.NotNil(t, Log)
	assert.Equal(t, zap.InfoLevel, Log.Level())
}

func TestRequestLogger(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	handler := RequestLogger(r)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "404 page not found", strings.TrimSpace(w.Body.String()))
}
func TestLoggingResponseWriter_Write(t *testing.T) {
	lw := &loggingResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		responseData: &responseData{
			status: 0,
			size:   0,
		},
	}
	b := []byte("Test response")
	n, err := lw.Write(b)
	assert.NoError(t, err)
	assert.Equal(t, len(b), n)
	assert.Equal(t, len(b), lw.responseData.size)
}

func TestLoggingResponseWriter_WriteHeader(t *testing.T) {
	lw := &loggingResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		responseData: &responseData{
			status: 0,
			size:   0,
		},
	}
	lw.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, lw.responseData.status)
}

func TestRequestLogger_LogsRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler := RequestLogger(r)

	start := time.Now()
	handler.ServeHTTP(w, req)
	elapsed := time.Since(start)

	// Check that the log contains the expected fields
	loggedEntry := Log.Check(zap.InfoLevel, "got incoming HTTP request")
	assert.NotNil(t, loggedEntry)
	loggedEntry.Write(zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("duration", elapsed.String()),
		zap.Int("status", http.StatusOK),
		zap.Int("size", len("Test response")),
	)
}
