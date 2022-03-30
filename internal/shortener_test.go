package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlShortenerService_ShortenURL(t *testing.T) {
	url := "https://www.google.com"
	repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, "short", short, "Expected short to be short, got %v", short)

	expanded, ok := repo.Expand(short)
	assert.Truef(t, ok, "Expected to find %s in the repo", short)
	assert.Equalf(t, url, expanded, "Expected short to be expanded to %s, got %s", url, expanded)

	repo.FailMode = true
	_, err = sut.ShortenURL(url)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ExpandShortURL(t *testing.T) {
	url := "https://www.google.com"
	repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	expanded, err := sut.ExpandShortURL(short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, url, expanded, "Expected %s to be expanded to %s, got %s", short, url, expanded)

	repo.FailMode = true
	_, err = sut.ExpandShortURL(short)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_DeleteShortURL(t *testing.T) {
	url := "https://www.google.com"
	repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, ok := repo.Expand(short)
	assert.Truef(t, ok, "Expected to find %s in the repo", short)

	err = sut.DeleteShortURL(short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, ok = repo.Expand(short)
	assert.Falsef(t, ok, "Expected to not find %s in the repo", short)
}

func setupService() (*redirectRepoFake, *UrlShortenerService) {
	hasher := func(_ string) string { return "short" }
	repo := NewInMemoryRedirectRepository()
	mock := &redirectRepoFake{repo: repo, FailMode: false}
	svc := NewUrlShortenerService(hasher, mock)
	return mock, &svc
}

type redirectRepoFake struct {
	repo     *InMemoryRedirectRepository
	FailMode bool
}

func (r redirectRepoFake) Save(short string, url string) error {
	if r.FailMode {
		return errors.New("mock error")
	}
	return r.repo.Save(short, url)
}

func (r redirectRepoFake) Expand(short string) (string, bool) {
	if r.FailMode {
		return "", false
	}
	return r.repo.Expand(short)
}

func (r redirectRepoFake) Delete(short string) error {
	if r.FailMode {
		return errors.New("mock error")
	}
	return r.repo.Delete(short)
}
