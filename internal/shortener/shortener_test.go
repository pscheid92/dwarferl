package shortener

import (
	"errors"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testUser = "00000000-0000-0000-0000-000000000000"
	testURL  = "https://www.google.com"
)

func TestUrlShortenerService_List(t *testing.T) {
	redirects, sut := setupService()

	list, err := sut.List(testUser)
	assert.NoErrorf(t, err, "list should not return error")
	assert.Emptyf(t, list, "list should return empty list")

	redirects.FailMode = true
	_, err = sut.List(testUser)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ShortenURL(t *testing.T) {
	repo, sut := setupService()

	redirect, err := sut.ShortenURL(testURL, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, "short", redirect.Short, "Expected short to be short, got %v", redirect.Short)

	expanded, err := repo.Expand(redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, testURL, expanded, "Expected redirect to be expanded to %s, got %s", testURL, expanded)

	repo.FailMode = true
	_, err = sut.ShortenURL(testURL, testUser)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ExpandShortURL(t *testing.T) {
	repo, sut := setupService()

	redirect, err := sut.ShortenURL(testURL, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	expanded, err := sut.ExpandShortURL(redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, testURL, expanded, "Expected %s to be expanded to %s, got %s", redirect.Short, testURL, expanded)

	repo.FailMode = true
	_, err = sut.ExpandShortURL(redirect.Short)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_DeleteShortURL(t *testing.T) {
	repo, sut := setupService()

	redirect, err := sut.ShortenURL(testURL, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, err = repo.Expand(redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	err = sut.DeleteShortURL(redirect.Short, testUser)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, err = repo.Expand(redirect.Short)
	assert.Errorf(t, err, "Expected error, got nil")
}

func setupService() (*redirectRepoFake, *UrlShortenerService) {
	hasher := func(_ string, _ string) string { return "short" }
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

func (r redirectRepoFake) List(user string) ([]internal.Redirect, error) {
	if r.FailMode {
		return nil, errors.New("fake error")
	}

	result := make([]internal.Redirect, 0, len(r.redirects))
	for _, redirect := range r.redirects {
		result = append(result, redirect)
	}
	return result, nil
}

func (r redirectRepoFake) Save(redirect internal.Redirect) error {
	if r.FailMode {
		return errors.New("fake error")
	}
	r.redirects[redirect.Short] = redirect
	return nil
}

func (r redirectRepoFake) Expand(short string) (string, error) {
	if r.FailMode {
		return "", errors.New("fake error")
	}

	url, ok := r.redirects[short]
	if !ok {
		return "", errors.New("not found")
	}
	return url.URL, nil
}

func (r redirectRepoFake) Delete(short string, userID string) error {
	if r.FailMode {
		return errors.New("fake error")
	}
	if _, ok := r.redirects[short]; !ok {
		return errors.New("not found")
	}
	delete(r.redirects, short)
	return nil
}
