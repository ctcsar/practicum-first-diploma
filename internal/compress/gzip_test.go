package compress

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCompressWriter(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	assert.NotNil(t, cw)
	assert.NotNil(t, cw.w)
	assert.NotNil(t, cw.zw)
}

func TestCompressWriter_Header(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	assert.NotNil(t, cw.Header())
	assert.Equal(t, w.Header(), cw.Header())
}

func TestCompressWriter_Write(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	data := []byte("Hello, World!")
	n, err := cw.Write(data)

	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	// Check that the data was compressed
	assert.NotEqual(t, data, w.Body.Bytes())
}

func TestCompressWriter_Write_Error(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	// Close the gzip writer to simulate an error
	cw.zw.Close()

	data := []byte("Hello, World!")
	_, err := cw.Write(data)

	assert.Error(t, err)
}

func TestCompressWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	cw.WriteHeader(http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCompressWriter_Close(t *testing.T) {
	w := httptest.NewRecorder()
	cw := newCompressWriter(w)

	err := cw.Close()

	assert.NoError(t, err)
}

func TestGzipMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			t.Errorf("failed to write to response writer: %v", err)
		}
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")

	GzipMiddleware(next).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
}

func TestGzipMiddleware_NoGzip(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			t.Errorf("failed to write to response writer: %v", err)
		}
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	GzipMiddleware(next).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Content-Encoding"))
}
