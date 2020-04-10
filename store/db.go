package store

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	f  string // db filename
	db *sql.DB
}

func Init(f string) *Store {
	db, err := sql.Open("sqlite3", f)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS games (token varchar not null primary key, desc text not null);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return &Store{f: f, db: db}
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Store) SyncGame(token string, game []byte) error {
	sqlStmt := fmt.Sprintf(`
	INSERT INTO games(token, desc) VALUES ('%s', '%s')
	ON CONFLICT(token) DO
	UPDATE SET desc='%s' WHERE token='%s';
	`, token, game, game, token)
	_, err := s.db.Exec(sqlStmt)
	return err
}

func (s *Store) LoadGames() (map[string][]byte, error) {
	games := make(map[string][]byte)
	rows, err := s.db.Query("SELECT token, desc FROM games;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var token string
		var game []byte
		err = rows.Scan(&token, &game)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		games[token] = game
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return games, nil
}
