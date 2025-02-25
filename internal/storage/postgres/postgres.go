package postgres

import (
	"awesomeProject/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(connStr string) (*PostgresStorage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open DB connection: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to ping DB: %w", op, err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) SaveURL(urlToSave string, shortUrl string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls(original_url, short_url) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(urlToSave, shortUrl).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *PostgresStorage) GetURL(shortUrl string) (string, error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.db.Prepare("SELECT original_url FROM urls WHERE short_url = $1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var resURL string
	err = stmt.QueryRow(shortUrl).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
