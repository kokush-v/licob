package db

import (
	"database/sql"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var (
	db *sql.DB
	rw sync.RWMutex
)

func Open(dsn string) error {
	var err error
	db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	return nil
}

func Up(dsn string) error {
	db, err := goose.OpenDBWithDriver("sqlite3", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Up(db, filepath.Join(".", "sql"))
}
