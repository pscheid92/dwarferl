package shortener

import (
	"context"
	"errors"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testUser = "00000000-0000-0000-0000-000000000000"
	testURL  = "https://www.google.com"
)

func TestUrlShortenerService_List(t *testing.T) {
	redirects, sut := setupService()

	list, err := sut.List(context.Background(), testUser)
	assert.NoErrorf(t, err, "list should not return error")
	assert.Emptyf(t, list, "list should return empty list")

	redirects.FailMode = true
	_, err = sut.List(context.Background(), testUser)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_GetRedirectByShort(t *testing.T) {
	redirects, sut := setupService()

	// invalid short
	_, err := sut.GetRedirectByShort(context.Background(), "invalid", testUser)
	assert.Error(t, err, "Expected error, got nil")

	// correct user and short
	redirect, err := sut.GetRedirectByShort(context.Background(), "short", testUser)
	assert.NoErrorf(t, err, "get redirect by short should not return error")
	assert.Equalf(t, testURL, redirect.URL, "get redirect by short should return correct redirect")
	assert.Equalf(t, testUser, redirect.UserID, "get redirect by short should return correct redirect")

	// false user correct short
	redirect, err = sut.GetRedirectByShort(context.Background(), "short", "nonexistent")
	assert.Errorf(t, err, "Expected error, got nil")

	// user and false short
	redirect, err = sut.GetRedirectByShort(context.Background(), "nonexistent", testUser)
	assert.Errorf(t, err, "Expected error, got nil")

	// error
	redirects.FailMode = true
	_, err = sut.GetRedirectByShort(context.Background(), "short", testUser)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ShortenURL(t *testing.T) {
	repo, sut := setupService()

	redirect, err := sut.ShortenURL(context.Background(), testURL, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, "short", redirect.Short, "Expected short to be short, got %v", redirect.Short)

	expanded, err := repo.Expand(context.Background(), redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, testURL, expanded, "Expected redirect to be expanded to %s, got %s", testURL, expanded)

	repo.FailMode = true
	_, err = sut.ShortenURL(context.Background(), testURL, testUser)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ExpandShortURL(t *testing.T) {
	repo, sut := setupService()

	_, err := sut.GetRedirectByShort(context.Background(), "invalid", testUser)
	assert.Error(t, err, "Expected error, got nil")

	redirect, err := sut.ShortenURL(context.Background(), testURL, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	expanded, err := sut.ExpandShortURL(context.Background(), redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, testURL, expanded, "Expected %s to be expanded to %s, got %s", redirect.Short, testURL, expanded)

	repo.FailMode = true
	_, err = sut.ExpandShortURL(context.Background(), redirect.Short)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_DeleteShortURL(t *testing.T) {
	repo, sut := setupService()

	redirect, err := sut.ShortenURL(context.Background(), testURL, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, err = repo.Expand(context.Background(), redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	err = sut.DeleteShortURL(context.Background(), redirect.Short, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, err = repo.Expand(context.Background(), redirect.Short)
	assert.Errorf(t, err, "Expected error, got nil")
}

func setupService() (*redirectRepoFake, *UrlShortenerService) {
	hasher := newHasherFake()
	redirects := newRedirectRepoFake()
	svc := NewUrlShortenerService(hasher, redirects)
	return redirects, &svc
}

type redirectRepoFake struct {
	redirects map[string]internal.Redirect
	FailMode  bool
}

func newRedirectRepoFake() *redirectRepoFake {
	return &redirectRepoFake{
		redirects: make(map[string]internal.Redirect),
		FailMode:  false,
	}
}

func (r redirectRepoFake) List(context.Context, string) ([]internal.Redirect, error) {
	if r.FailMode {
		return nil, errors.New("fake error")
	}

	result := make([]internal.Redirect, 0, len(r.redirects))
	for _, redirect := range r.redirects {
		result = append(result, redirect)
	}
	return result, nil
}

func (r redirectRepoFake) GetRedirectByShort(_ context.Context, short string, userID string) (internal.Redirect, error) {
	if r.FailMode {
		return internal.Redirect{}, errors.New("fake error")
	}

	if short != "short" || userID != testUser {
		return internal.Redirect{}, errors.New("not found")
	}

	result := internal.Redirect{
		Short:     "short",
		URL:       testURL,
		UserID:    testUser,
		CreatedAt: time.Now().Add(-time.Hour),
	}
	return result, nil
}

func (r redirectRepoFake) Save(_ context.Context, redirect internal.Redirect) error {
	if r.FailMode {
		return errors.New("fake error")
	}
	r.redirects[redirect.Short] = redirect
	return nil
}

func (r redirectRepoFake) Expand(_ context.Context, short string) (string, error) {
	if r.FailMode {
		return "", errors.New("fake error")
	}

	url, ok := r.redirects[short]
	if !ok {
		return "", errors.New("not found")
	}
	return url.URL, nil
}

func (r redirectRepoFake) Delete(_ context.Context, short string, _ string) error {
	if r.FailMode {
		return errors.New("fake error")
	}
	if _, ok := r.redirects[short]; !ok {
		return errors.New("not found")
	}
	delete(r.redirects, short)
	return nil
}

type hasherFake struct{}

func newHasherFake() *hasherFake {
	return &hasherFake{}
}

func (h hasherFake) Hash(string, string) string {
	return "short"
}

func (h hasherFake) Validate(short string) bool {
	return short == "short"
}
