package internal

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupStaticHasher(fixed string) Hasher {
	return func(string) string {
		return fixed
	}
}

func TestCreateHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", CreateHealthHandler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code to be 200, got: %d", w.Code)
	}
}

func TestCreateGetHandler(t *testing.T) {
	url := "https://www.google.com"
	repo := NewInMemoryRedirectRepository()
	svc := NewUrlShortenerService(setupStaticHasher("short"), repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/:short", CreateGetHandler(svc))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code to be 404, got: %d", w.Code)
	}

	short, err := svc.ShortenURL(url)
	if err != nil {
		t.Error("Expected no error, got: ", err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/"+short, nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code to be 302, got: %d", w.Code)
	}
	if w.Header().Get("Location") != url {
		t.Errorf("Expected Location header to be %s, got: %s", url, w.Header().Get("Location"))
	}
}

func TestCreatePostHandler(t *testing.T) {
	repo := NewInMemoryRedirectRepository()
	mock := &repoMock{repo: repo, FailMode: false}
	svc := NewUrlShortenerService(setupStaticHasher("short"), mock)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/", CreatePostHandler(svc))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code to be 400, got: %d", w.Code)
	}

	w = httptest.NewRecorder()
	body := bytes.NewBuffer([]byte(`{}`))
	req, _ = http.NewRequest("POST", "/", body)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code to be 400, got: %d", w.Code)
	}

	w = httptest.NewRecorder()
	body = bytes.NewBuffer([]byte(`{"url": "https://www.google.com"}`))
	req, _ = http.NewRequest("POST", "/", body)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code to be 201, got: %d", w.Code)
	}

	w = httptest.NewRecorder()
	body = bytes.NewBuffer([]byte(`{"url": "https://www.google.com"}`))
	req, _ = http.NewRequest("POST", "/", body)
	mock.FailMode = true
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code to be 500, got: %d", w.Code)
	}
}

func TestCreateDeleteHandler(t *testing.T) {
	url := "https://www.google.com"
	repo := NewInMemoryRedirectRepository()
	svc := NewUrlShortenerService(setupStaticHasher("short"), repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/:short", CreateDeleteHandler(svc))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/nonexistent", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code to be 404, got: %d", w.Code)
	}

	_, err := svc.ShortenURL(url)
	if err != nil {
		t.Error("Expected no error, got: ", err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/short", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code to be 200, got: %d", w.Code)
	}
}
