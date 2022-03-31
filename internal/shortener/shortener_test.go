package shortener

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlShortenerService_List(t *testing.T) {
	users, redirects, sut := setupService()

	list, err := sut.List("user1")
	assert.NoErrorf(t, err, "list should not return error")
	assert.Emptyf(t, list, "list should return empty list")

	redirects.FailMode = true
	_, err = sut.List("user1")
	assert.Errorf(t, err, "Expected error, got nil")

	users.FailMode = true
	_, err = sut.List("user1")
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ShortenURL(t *testing.T) {
	url := "https://www.google.com"
	_, repo, sut := setupService()

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
	_, repo, sut := setupService()

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
	_, repo, sut := setupService()

	short, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, ok := repo.Expand(short)
	assert.Truef(t, ok, "Expected to find %s in the repo", short)

	err = sut.DeleteShortURL(short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, ok = repo.Expand(short)
	assert.Falsef(t, ok, "Expected to not find %s in the repo", short)
}

func setupService() (*usersRepoFake, *redirectRepoFake, *UrlShortenerService) {
	hasher := func(_ internal.User, _ string) string { return "short" }

	redirects := newRedirectRepoFake()
	users := newUsersRepoFake()

	svc := NewUrlShortenerService(hasher, redirects, users)
	return users, redirects, &svc
}

type redirectRepoFake struct {
	redirects map[string]string
	FailMode  bool
}

func newRedirectRepoFake() *redirectRepoFake {
	return &redirectRepoFake{
		redirects: make(map[string]string),
		FailMode:  false,
	}
}

func (r redirectRepoFake) List(_ internal.User) (map[string]string, error) {
	if r.FailMode {
		return nil, errors.New("fake error")
	}
	return r.redirects, nil
}

func (r redirectRepoFake) Save(short string, url string) error {
	if r.FailMode {
		return errors.New("fake error")
	}
	r.redirects[short] = url
	return nil
}

func (r redirectRepoFake) Expand(short string) (string, bool) {
	if r.FailMode {
		return "", false
	}
	url, ok := r.redirects[short]
	return url, ok
}

func (r redirectRepoFake) Delete(short string) error {
	if r.FailMode {
		return errors.New("fake error")
	}
	if _, ok := r.redirects[short]; !ok {
		return errors.New("not found")
	}
	delete(r.redirects, short)
	return nil
}

type usersRepoFake struct {
	FailMode bool
}

func newUsersRepoFake() *usersRepoFake {
	return &usersRepoFake{FailMode: false}
}

func (u usersRepoFake) Get(string) (internal.User, error) {
	if u.FailMode {
		return internal.User{}, errors.New("fake error")
	}

	user := internal.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Email: "example@example.com",
	}
	return user, nil
}
