package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHealthHandler(t *testing.T) {
	_, router := setupTest()
	w := executeCall(router, "GET", "/health", "")
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200")
}

func TestCreateGetHandler(t *testing.T) {
	_, router := setupTest()

	w := executeCall(router, "GET", "/nonexistent", "")
	assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)

	w = executeCall(router, "GET", "/"+testShort, "")

	location := w.Header().Get("Location")
	cacheControl := w.Header().Get("Cache-Control")
	referrerPolicy := w.Header().Get("Referrer-Policy")

	assert.Equalf(t, http.StatusMovedPermanently, w.Code, "Expected status code to be 301, got %d", w.Code)
	assert.Equalf(t, testURL, location, "Expected location header to be %s, got %s", testURL, location)
	assert.Containsf(t, cacheControl, "private", "Expected cache-control header to contain private, got %s", cacheControl)
	assert.Containsf(t, cacheControl, "max-age", "Expected cache-control header to contain max-age, got %s", cacheControl)
	assert.Equalf(t, "unsafe-url", referrerPolicy, "Expected referrer-policy header to be unsafe-url, got %s", referrerPolicy)
}

func TestCreateListHandler(t *testing.T) {
	svc, router := setupTest()

	w := executeCall(router, "GET", "/", "")
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)

	svc.FailMode = true
	w = executeCall(router, "GET", "/", "")
	assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
}

func TestCreatePostHandler(t *testing.T) {
	svc, router := setupTest()

	w := executeCall(router, "POST", "/", "")
	assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected status code to be 400, got %d", w.Code)

	w = executeCall(router, "POST", "/", `{}`)
	assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected status code to be 400, got %d", w.Code)

	w = executeCall(router, "POST", "/", `{"url": "https://www.google.com"}`)
	assert.Equalf(t, http.StatusCreated, w.Code, "Expected status code to be 201, got %d", w.Code)

	svc.FailMode = true
	w = executeCall(router, "POST", "/", `{"url": "https://www.google.com"}`)
	assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
}

func TestCreateDeleteHandler(t *testing.T) {
	_, router := setupTest()

	w := executeCall(router, "DELETE", "/nonexistent", "")
	assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)

	w = executeCall(router, "DELETE", "/"+testShort, "")
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
}

func setupTest() (*urlShortenerServiceFake, *gin.Engine) {
	svc := newShortenerServiceFake()

	gin.SetMode(gin.TestMode)
	accounts := gin.Accounts{"test": "test"}
	router := SetupRoutes(gin.New(), "/", svc, accounts)

	return svc, router
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

const testUser = "user1"
const testShort = "short"
const testURL = "https://www.google.com"

type urlShortenerServiceFake struct {
	FailMode bool
}

func newShortenerServiceFake() *urlShortenerServiceFake {
	return &urlShortenerServiceFake{FailMode: false}
}

func (s urlShortenerServiceFake) List(userID string) (map[string]string, error) {
	if s.FailMode {
		return nil, errors.New("fake error")
	}

	if userID != testUser {
		return nil, errors.New("user not found")
	}

	return map[string]string{testShort: testURL}, nil
}

func (s urlShortenerServiceFake) ShortenURL(url string) (string, error) {
	if s.FailMode {
		return "", errors.New("fake error")
	}

	if url != testURL {
		return "", errors.New("not found")
	}

	return testShort, nil
}

func (s urlShortenerServiceFake) ExpandShortURL(short string) (string, error) {
	if s.FailMode {
		return "", errors.New("fake error")
	}

	if short != testShort {
		return "", errors.New("not found")
	}

	return testURL, nil
}

func (s urlShortenerServiceFake) DeleteShortURL(short string) error {
	if s.FailMode {
		return errors.New("fake error")
	}

	if short != testShort {
		return errors.New("not found")
	}

	return nil
}
