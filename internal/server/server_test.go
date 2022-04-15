package server

import (
	"context"
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
	srv, _, _ := setupTestServer()
	w := srv.call("GET", "/health", "", nil)
	assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200")
}

func TestHandleRedirect(t *testing.T) {
	srv, _, _ := setupTestServer()

	t.Run("non-existent redirect", func(t *testing.T) {
		w := srv.call("GET", "/nonexistent", "", nil)
		assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 404, got %d", w.Code)
	})

	t.Run("invalid characters are rejected", func(t *testing.T) {
		w := srv.call("GET", "/falséy", "", nil)
		assert.Equalf(t, http.StatusNotFound, w.Code, "Expected status code to be 400, got %d", w.Code)
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
	srv, cookies, _ := setupTestServer()

	t.Run("login page serves successfully if not logged in", func(t *testing.T) {
		w := srv.call("GET", "/login", "", nil)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d %v", w.Code, w)
	})

	t.Run("login page redirects to homepage if already logged in", func(t *testing.T) {
		w := srv.call("GET", "/login", "", cookies)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d %v", w.Code, w)
	})
}

func TestHandleIndexPage(t *testing.T) {
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
	srv, cookies, _ := setupTestServer()

	t.Run("deletion page demands login", func(t *testing.T) {
		w := srv.call("GET", "/delete/"+testShort, "url="+testURL, nil)
		assert.Equalf(t, http.StatusFound, w.Code, "Expected status code to be 302, got %d", w.Code)
	})

	t.Run("deletion of nonexistent short fails", func(t *testing.T) {
		w := srv.call("GET", "/delete/nonexistent", "", cookies)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected status code to be 500, got %d", w.Code)
	})

	t.Run("deletion page served successfully", func(t *testing.T) {
		w := srv.call("GET", "/delete/"+testShort, "", cookies)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected status code to be 200, got %d", w.Code)
	})
}

func TestHandlePostDeletionPage(t *testing.T) {
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
		TemplatePath:    "../../templates",
	}

	shortener := &urlShortenerServiceFake{}
	users := &usersServiceFake{}
	store := cookie.NewStore([]byte(c.SessionSecret))

	gin.SetMode(gin.TestMode)
	svr := New(c, store, shortener, users)
	svr.InitRoutes()

	cookies := svr.autologin()
	return svr, cookies, shortener
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

type urlShortenerServiceFake struct {
	FailMode bool
}

func (s urlShortenerServiceFake) List(_ context.Context, userID string) ([]internal.Redirect, error) {
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

func (s urlShortenerServiceFake) GetRedirectByShort(_ context.Context, short string, userID string) (internal.Redirect, error) {
	if s.FailMode {
		return internal.Redirect{}, errors.New("fake error")
	}

	if short != "short" || userID != testUser {
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

func (s urlShortenerServiceFake) ShortenURL(_ context.Context, url string, _ string) (internal.Redirect, error) {
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

func (s urlShortenerServiceFake) ExpandShortURL(_ context.Context, short string) (string, error) {
	if s.FailMode {
		return "", errors.New("fake error")
	}

	if short != testShort {
		return "", errors.New("not found")
	}

	return testURL, nil
}

func (s urlShortenerServiceFake) DeleteShortURL(_ context.Context, short string, _ string) error {
	if s.FailMode {
		return errors.New("fake error")
	}

	if short != testShort {
		return errors.New("not found")
	}

	return nil
}

type usersServiceFake struct {
	FailMode bool
}

func (u usersServiceFake) CreateWithGoogleID(context.Context, string, string) (internal.User, error) {
	// TODO: find a way to test goth/gothic
	// Not implemented, since we cannot test goth/gothic here
	panic("create with google id - implement me")
}

func (u usersServiceFake) GetOrCreateByGoogle(context.Context, string, string) (internal.User, error) {
	// TODO: find a way to test goth/gothic
	// Not implemented, since we cannot test goth/gothic here
	panic("get or create by google - implement me")
}
