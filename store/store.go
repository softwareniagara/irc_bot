package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func Open(filename string) (*Store, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	s := &Store{db}
	if err := s.CreateTables(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) CreateTables() error {
	if _, err := s.db.Exec(createReminderSQL); err != nil {
		return err
	}
	if _, err := s.db.Exec(createUserSQL); err != nil {
		return err
	}
	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
