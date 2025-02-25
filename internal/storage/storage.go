package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Storage interface {
	SaveURL(urlToSave string, shortUrl string) (int64, error)
	GetURL(shortUrl string) (string, error)
}
