package internal

import (
	"errors"
	"testing"
)

func TestUrlShortenerService_ShortenURL(t *testing.T) {
	url := "https://www.google.com"
	repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if short != url {
		t.Errorf("Expected short to be %s, got %s", short, url)
	}

	if expanded, ok := repo.Expand(short); !ok {
		t.Errorf("Expected to find %s in the repo", short)
	} else if expanded != url {
		t.Errorf("Expected %s to be expanded to %s, got %s", short, url, expanded)
	}

	repo.FailMode = true
	if _, err = sut.ShortenURL(url); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUrlShortenerService_ExpandShortURL(t *testing.T) {
	url := "https://www.google.com"
	repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expanded, err := sut.ExpandShortURL(short)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if expanded != url {
		t.Errorf("Expected %s to be expanded to %s, got %s", short, url, expanded)
	}

	repo.FailMode = true
	if _, err = sut.ExpandShortURL(short); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUrlShortenerService_DeleteShortURL(t *testing.T) {
	url := "https://www.google.com"
	repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, ok := repo.Expand(short)
	if !ok {
		t.Errorf("Expected to find %s in the repo", short)
	}

	err = sut.DeleteShortURL(short)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if _, ok = repo.Expand(short); ok {
		t.Errorf("Expected to not find %s in the repo", short)
	}
}

func noOpHasher(input string) string {
	return input
}

func setupService() (*repoMock, *UrlShortenerService) {
	repo := NewInMemoryRedirectRepository()
	mock := &repoMock{repo: repo, FailMode: false}
	svc := NewUrlShortenerService(noOpHasher, mock)
	return mock, &svc
}

type repoMock struct {
	repo     *InMemoryRedirectRepository
	FailMode bool
}

func (r repoMock) Save(short string, url string) error {
	if r.FailMode {
		return errors.New("mock error")
	}
	return r.repo.Save(short, url)
}

func (r repoMock) Expand(short string) (string, bool) {
	if r.FailMode {
		return "", false
	}
	return r.repo.Expand(short)
}

func (r repoMock) Delete(short string) error {
	if r.FailMode {
		return errors.New("mock error")
	}
	return r.repo.Delete(short)
}
