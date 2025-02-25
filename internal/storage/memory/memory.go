package memory

import (
	"errors"
	"log"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	links map[string]string
}

func New() *Store {
	return &Store{
		links: make(map[string]string),
	}
}

func (s *Store) SaveURL(urlToSave string, shortUrl string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Attempting to save URL: %s with short URL: %s", urlToSave, shortUrl)

	if _, exists := s.links[shortUrl]; exists {
		return 0, errors.New("short URL already exists")
	}

	s.links[shortUrl] = urlToSave
	return 1, nil
}

func (s *Store) GetURL(shortUrl string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("Attempting to retrieve original URL for short URL: %s", shortUrl)

	originalURL, exists := s.links[shortUrl]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return originalURL, nil
}
