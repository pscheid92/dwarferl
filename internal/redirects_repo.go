package internal

import "errors"

type RedirectRepository interface {
	Save(short string, url string) error
	Expand(short string) (string, bool)
	Delete(short string) error
}

type InMemoryRedirectRepository struct {
	redirects map[string]string
}

func NewInMemoryRedirectRepository() *InMemoryRedirectRepository {
	return &InMemoryRedirectRepository{
		redirects: make(map[string]string),
	}
}

func (i InMemoryRedirectRepository) Save(short string, url string) error {
	i.redirects[short] = url
	return nil
}

func (i InMemoryRedirectRepository) Expand(short string) (string, bool) {
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
