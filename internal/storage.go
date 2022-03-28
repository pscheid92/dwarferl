package internal

type Storage interface {
	SaveRedirect(short string, url string) error
	ExpandRedirect(short string) (string, bool)
	DeleteRedirect(short string) error
}

type InMemoryStorage struct {
	redirects map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		redirects: make(map[string]string),
	}
}

func (i InMemoryStorage) SaveRedirect(short string, url string) error {
	i.redirects[short] = url
	return nil
}

func (i InMemoryStorage) ExpandRedirect(short string) (string, bool) {
	url, ok := i.redirects[short]
	return url, ok
}

func (i InMemoryStorage) DeleteRedirect(short string) error {
	delete(i.redirects, short)
	return nil
}
