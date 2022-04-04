package handler

import (
	"errors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateHealthHandler(t *testing.T) {
	_, router, _ := setupTest()
	w := executeCall(router, "GET", "/health", "", nil)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200")
}

func TestCreateGetHandler(t *testing.T) {
	_, router, _ := setupTest()

	w := executeCall(router, "GET", "/nonexistent", "", nil)
	assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)

	w = executeCall(router, "GET", "/"+testShort, "", nil)

	location := w.Header().Get("Location")
	cacheControl := w.Header().Get("Cache-Control")
	referrerPolicy := w.Header().Get("Referrer-Policy")

	assert.Equalf(t, http.StatusMovedPermanently, w.Code, "Expected status code to be 301, got %d", w.Code)
	assert.Equalf(t, testURL, location, "Expected location header to be %s, got %s", testURL, location)
	assert.Containsf(t, cacheControl, "private", "Expected cache-control header to contain private, got %s", cacheControl)
	assert.Containsf(t, cacheControl, "max-age", "Expected cache-control header to contain max-age, got %s", cacheControl)
	assert.Equalf(t, "unsafe-url", referrerPolicy, "Expected referrer-policy header to be unsafe-url, got %s", referrerPolicy)
}

func TestIndexPage(t *testing.T) {
	svc, router, cookie := setupTest()

	w := executeCall(router, "GET", "/", "", cookie)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)

	svc.FailMode = true
	w = executeCall(router, "GET", "/", "", cookie)
	assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
}

func TestServeCreationPage(t *testing.T) {
	_, router, cookie := setupTest()
	w := executeCall(router, "GET", "/create", "", cookie)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
}

func TestHandleCreationPage(t *testing.T) {
	svc, router, cookie := setupTest()

	w := executeCall(router, "POST", "/create", "url="+testURL, cookie)
	assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)

	svc.FailMode = true
	w = executeCall(router, "POST", "/create", "url="+testURL, cookie)
	assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
}

func TestServeDeletionPage(t *testing.T) {
	_, router, cookie := setupTest()
	w := executeCall(router, "GET", "/delete/short", "", cookie)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
}

func TestCreateDeleteHandler(t *testing.T) {
	_, router, cookie := setupTest()

	w := executeCall(router, "POST", "/delete/nonexistent", "", cookie)
	assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)

	w = executeCall(router, "POST", "/delete/"+testShort, "", cookie)
	assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
}

func setupTest() (*urlShortenerServiceFake, *gin.Engine, *http.Cookie) {
	svc := newShortenerServiceFake()

	cookies := sessions.NewCookieStore([]byte("test_secret"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	router = SetupRoutes(router, "/", svc, cookies)

	router.GET("/autologin", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user_id", testUser)

		if err := session.Save(); err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.Status(200)
	})

	w := executeCall(router, "GET", "/autologin", "", nil)
	cookie := w.Result().Cookies()[0]

	return svc, router, cookie
}

func executeCall(router *gin.Engine, method string, url string, body string, cookies *http.Cookie) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, reader)

	if cookies != nil {
		req.AddCookie(cookies)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	return w
}

const testUser = "00000000-0000-0000-0000-000000000000"
const testShort = "short"
const testURL = "https://www.google.com"

type urlShortenerServiceFake struct {
	FailMode bool
}

func newShortenerServiceFake() *urlShortenerServiceFake {
	return &urlShortenerServiceFake{FailMode: false}
}

func (s urlShortenerServiceFake) List(userID string) ([]internal.Redirect, error) {
	if s.FailMode {
		return nil, errors.New("fake error")
	}

	if userID != testUser {
		return nil, errors.New("user not found")
	}

	redirect := internal.Redirect{
		Short:     testShort,
		URL:       testURL,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	return []internal.Redirect{redirect}, nil
}

func (s urlShortenerServiceFake) ShortenURL(url string) (internal.Redirect, error) {
	if s.FailMode {
		return internal.Redirect{}, errors.New("fake error")
	}

	if url != testURL {
		return internal.Redirect{}, errors.New("not found")
	}

	redirect := internal.Redirect{
		Short:     testShort,
		URL:       testURL,
		UserID:    testUser,
		CreatedAt: time.Now(),
	}

	return redirect, nil
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
