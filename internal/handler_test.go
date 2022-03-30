package internal

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest() (*redirectRepoFake, *UrlShortenerService, *gin.Engine) {
	staticHasher := func(string) string { return "short" }

	repo := NewInMemoryRedirectRepository()
	mock := redirectRepoFake{repo: repo, FailMode: false}
	svc := NewUrlShortenerService(staticHasher, &mock)

	gin.SetMode(gin.TestMode)
	accounts := gin.Accounts{"test": "test"}
	router := SetupRoutes(gin.New(), "/", svc, accounts)

	return &mock, &svc, router
}

func executeCall(router *gin.Engine, method string, url string, body string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != "" {
		reader = bytes.NewReader([]byte(body))
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, reader)
	req.SetBasicAuth("test", "test")
	router.ServeHTTP(w, req)
	return w
}

func TestCreateHealthHandler(t *testing.T) {
	_, _, router := setupTest()
	w := executeCall(router, "GET", "/health", "")
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200")
}

func TestCreateGetHandler(t *testing.T) {
	url := "https://www.google.com"
	_, svc, router := setupTest()

	w := executeCall(router, "GET", "/nonexistent", "")
	assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)

	short, err := svc.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	w = executeCall(router, "GET", "/"+short, "")

	location := w.Header().Get("Location")
	cacheControl := w.Header().Get("Cache-Control")
	referrerPolicy := w.Header().Get("Referrer-Policy")

	assert.Equalf(t, http.StatusMovedPermanently, w.Code, "Expected status code to be 301, got %d", w.Code)
	assert.Equalf(t, url, location, "Expected location header to be %s, got %s", url, location)
	assert.Containsf(t, cacheControl, "private", "Expected cache-control header to contain private, got %s", cacheControl)
	assert.Containsf(t, cacheControl, "max-age", "Expected cache-control header to contain max-age, got %s", cacheControl)
	assert.Equalf(t, "unsafe-url", referrerPolicy, "Expected referrer-policy header to be unsafe-url, got %s", referrerPolicy)
}

func TestCreatePostHandler(t *testing.T) {
	repo, _, router := setupTest()

	w := executeCall(router, "POST", "/", "")
	assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected status code to be 400, got %d", w.Code)

	w = executeCall(router, "POST", "/", `{}`)
	assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected status code to be 400, got %d", w.Code)

	w = executeCall(router, "POST", "/", `{"url": "https://www.google.com"}`)
	assert.Equalf(t, http.StatusCreated, w.Code, "Expected status code to be 201, got %d", w.Code)

	repo.FailMode = true
	w = executeCall(router, "POST", "/", `{"url": "https://www.google.com"}`)
	assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
}

func TestCreateDeleteHandler(t *testing.T) {
	url := "https://www.google.com"
	_, svc, router := setupTest()

	w := executeCall(router, "DELETE", "/nonexistent", "")
	assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)

	_, err := svc.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	w = executeCall(router, "DELETE", "/short", "")
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
}
