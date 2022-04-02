package repository

import (
	"errors"
	"github.com/pscheid92/dwarferl/internal"
)

type InMemoryRedirectRepository struct {
	redirects map[string]internal.Redirect
}

func NewInMemoryRedirectRepository() *InMemoryRedirectRepository {
	return &InMemoryRedirectRepository{
		redirects: make(map[string]internal.Redirect),
	}
}

func (i InMemoryRedirectRepository) List(user internal.User) ([]internal.Redirect, error) {
	result := make([]internal.Redirect, 0, len(i.redirects))

	for _, redirect := range i.redirects {
		result = append(result, redirect)
	}

	return result, nil
}

func (i InMemoryRedirectRepository) Save(redirect internal.Redirect) error {
	i.redirects[redirect.Short] = redirect
	return nil
}

func (i InMemoryRedirectRepository) Expand(short string) (internal.Redirect, bool) {
	url, ok := i.redirects[short]
	return url, ok
}

func (i InMemoryRedirectRepository) Delete(short string) error {
	if _, ok := i.redirects[short]; !ok {
		return errors.New("not found")
	}
	delete(i.redirects, short)
	return nil
}
