package store

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	f  string // db filename
	db *sql.DB
}

func Init(f string) *Store {
	return &Store{f: f}
}

func (s *Store) Connect() {
	db, err := sql.Open("sqlite3", s.f)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
