package server

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/pscheid92/dwarferl/internal/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	testUser  = "00000000-0000-0000-0000-000000000000"
	testShort = "short"
	testURL   = "https://www.google.com"
)

func TestHandleHealth(t *testing.T) {
	t.Parallel()
	srv, _, _ := setupTestServer()
	w := srv.call("GET", "/health", "", nil)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200")
}

func TestHandleRedirect(t *testing.T) {
	t.Parallel()
	srv, _, _ := setupTestServer()

	t.Run("non-existent redirect", func(t *testing.T) {
		w := srv.call("GET", "/nonexistent", "", nil)
		assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)
	})

	t.Run("redirect successfully", func(t *testing.T) {
		w := srv.call("GET", "/"+testShort, "", nil)

		location := w.Header().Get("Location")
		cacheControl := w.Header().Get("Cache-Control")
		referrerPolicy := w.Header().Get("Referrer-Policy")

		assert.Equalf(t, http.StatusMovedPermanently, w.Code, "Expected status code to be 301, got %d", w.Code)
		assert.Equalf(t, testURL, location, "Expected location header to be %s, got %s", testURL, location)
		assert.Containsf(t, cacheControl, "private", "Expected cache-control header to contain private, got %s", cacheControl)
		assert.Containsf(t, cacheControl, "max-age", "Expected cache-control header to contain max-age, got %s", cacheControl)
		assert.Equalf(t, "unsafe-url", referrerPolicy, "Expected referrer-policy header to be unsafe-url, got %s", referrerPolicy)
	})
}

func TestHandleGetLoginPage(t *testing.T) {
	t.Parallel()
	srv, _, _ := setupTestServer()
	w := srv.call("GET", "/login", "", nil)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d %v", w.Code, w)
}

func TestHandlePostLoginPage(t *testing.T) {
	t.Parallel()
	srv, _, _ := setupTestServer()

	t.Run("login with wrong credentials", func(t *testing.T) {
		w := srv.call("POST", "/login", "username=wrong&password=wrong", nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
		assert.NotContainsf(t, extractCookieNames(t, w), "dwarferl_session", "Expected session cookie to be missing")
	})

	t.Run("login successfully", func(t *testing.T) {
		w := srv.call("POST", "/login", "username=admin&password=admin", nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
		assert.Containsf(t, extractCookieNames(t, w), "dwarferl_session", "Expected session cookie to be present")
	})
}

func TestHandleLogoutPage(t *testing.T) {
	t.Parallel()
	srv, cookies, _ := setupTestServer()

	t.Run("logout without being logged in", func(t *testing.T) {
		w := srv.call("GET", "/logout", "", nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expectet status code 302, got %d", w.Code)
		assert.NotContainsf(t, extractCookieNames(t, w), "dwarferl_session", "Expected session cookie to be missing")
	})

	t.Run("logout successfully", func(t *testing.T) {
		w := srv.call("GET", "/logout", "", cookies)
		assert.Equalf(t, http.StatusFound, w.Code, "Expectet status code 302, got %d", w.Code)

		newCookies := w.Result().Cookies()
		w = srv.call("GET", "/", "", newCookies)
		assert.Equalf(t, http.StatusFound, w.Code, "Expectet status code 302, got %d", w.Code)
		assert.Equalf(t, "/login", w.Header().Get("Location"), "Expected redirect to login page")
	})
}

func TestHandleIndexPage(t *testing.T) {
	t.Parallel()
	srv, cookies, shortener := setupTestServer()

	t.Run("index page demands login", func(t *testing.T) {
		w := srv.call("GET", "/", "", nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("index page is shown successfully", func(t *testing.T) {
		w := srv.call("GET", "/", "", cookies)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d %v", w.Code, w)
	})

	t.Run("index page show dead-end error", func(t *testing.T) {
		shortener.FailMode = true
		w := srv.call("GET", "/", "", cookies)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
	})
}

func TestHandleGetCreationPage(t *testing.T) {
	t.Parallel()
	srv, cookies, _ := setupTestServer()

	t.Run("creation page demands login", func(t *testing.T) {
		w := srv.call("GET", "/create", "", nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("creation page served successfully", func(t *testing.T) {
		w := srv.call("GET", "/create", "", cookies)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
	})

}

func TestHandlePostCreationPage(t *testing.T) {
	t.Parallel()
	srv, cookies, shortener := setupTestServer()

	t.Run("creation post demands login", func(t *testing.T) {
		w := srv.call("POST", "/create", "url="+testURL, nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("creation post processed successfully", func(t *testing.T) {
		w := srv.call("POST", "/create", "url="+testURL, cookies)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("creation post shows dead-end error", func(t *testing.T) {
		shortener.FailMode = true
		w := srv.call("POST", "/create", "url="+testURL, cookies)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
	})
}

func TestHandleGetDeletionPage(t *testing.T) {
	t.Parallel()
	srv, cookies, _ := setupTestServer()

	t.Run("deletion page demands login", func(t *testing.T) {
		w := srv.call("GET", "/delete/"+testShort, "url="+testURL, nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("deletion page served successfully", func(t *testing.T) {
		w := srv.call("GET", "/delete/"+testShort, "", cookies)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
	})
}

func TestHandlePostDeletionPage(t *testing.T) {
	t.Parallel()
	srv, cookies, _ := setupTestServer()

	t.Run("deletion post demands login", func(t *testing.T) {
		w := srv.call("GET", "/delete/"+testShort, "url="+testURL, nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("deletion post for non-existing redirect show not found", func(t *testing.T) {
		w := srv.call("POST", "/delete/nonexistent", "", cookies)
		assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)
	})

	t.Run("deletion post successfully deletes", func(t *testing.T) {
		w := srv.call("POST", "/delete/"+testShort, "", cookies)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})
}

func setupTestServer() (*Server, []*http.Cookie, *urlShortenerServiceFake) {
	c := config.Configuration{
		ForwardedPrefix: "/",
		SessionSecret:   "test_secret",
		TemplatePath:    "../../templates",
	}

	svc := newShortenerServiceFake()
	store := cookie.NewStore([]byte(c.SessionSecret))

	gin.SetMode(gin.TestMode)
	svr := New(c, store, svc)
	svr.InitRoutes()

	cookies := svr.autologin()
	return svr, cookies, svc
}

func (s *Server) autologin() []*http.Cookie {
	// setup route for automatic login in tests
	s.POST("/test_login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user_id", testUser)
		if err := session.Save(); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(200)
	})

	// login and extract valid session cookies
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/test_login", nil)
	s.ServeHTTP(w, r)
	return w.Result().Cookies()
}

func (s *Server) call(method string, url string, body string, cookies []*http.Cookie) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, url, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// add session cookies to call
	for _, c := range cookies {
		r.AddCookie(c)
	}

	s.ServeHTTP(w, r)
	return w
}

func extractCookieNames(t *testing.T, w *httptest.ResponseRecorder) []string {
	t.Helper()

	names := make([]string, 0)
	for _, c := range w.Result().Cookies() {
		names = append(names, c.Name)
	}
	return names
}

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
