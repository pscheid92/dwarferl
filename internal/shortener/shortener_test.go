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

	redirect, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, "short", redirect.Short, "Expected short to be short, got %v", redirect.Short)

	expanded, ok := repo.Expand(redirect.Short)
	assert.Truef(t, ok, "Expected to find %s in the repo", redirect.Short)
	assert.Equalf(t, url, expanded.URL, "Expected redirect to be expanded to %s, got %s", url, expanded.URL)

	repo.FailMode = true
	_, err = sut.ShortenURL(url)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_ExpandShortURL(t *testing.T) {
	url := "https://www.google.com"
	_, repo, sut := setupService()

	redirect, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	expanded, err := sut.ExpandShortURL(redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, url, expanded.URL, "Expected %s to be expanded to %s, got %s", redirect.Short, url, expanded)

	repo.FailMode = true
	_, err = sut.ExpandShortURL(redirect.Short)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestUrlShortenerService_DeleteShortURL(t *testing.T) {
	url := "https://www.google.com"
	_, repo, sut := setupService()

	redirect, err := sut.ShortenURL(url)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, ok := repo.Expand(redirect.Short)
	assert.Truef(t, ok, "Expected to find %s in the repo", redirect.Short)

	err = sut.DeleteShortURL(redirect.Short)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)

	_, ok = repo.Expand(redirect.Short)
	assert.Falsef(t, ok, "Expected to not find %s in the repo", redirect.Short)
}

func setupService() (*usersRepoFake, *redirectRepoFake, *UrlShortenerService) {
	hasher := func(_ internal.User, _ string) string { return "short" }

	redirects := newRedirectRepoFake()
	users := newUsersRepoFake()

	svc := NewUrlShortenerService(hasher, redirects, users)
	return users, redirects, &svc
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

func (r redirectRepoFake) List(user internal.User) ([]internal.Redirect, error) {
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

func (r redirectRepoFake) Expand(short string) (internal.Redirect, bool) {
	if r.FailMode {
		return internal.Redirect{}, false
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
