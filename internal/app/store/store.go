package store

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Store
type Store struct {
	config  *Config
	db      *sql.DB
	innRepo *InnRepo
}

// New ...
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open ...
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

//Close ...
func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Store) Inn() *InnRepo {
	if s.innRepo != nil {
		return s.innRepo
	}

	s.innRepo = &InnRepo{
		store: s,
	}

	return s.innRepo
}
