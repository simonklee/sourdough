package main

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/simonklee/sourdough/query"
)

//go:embed schema.sql
var ddl string

func defaultDBPath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "sourdough/sourdough.sqlite")
}

func InitStore(ctx context.Context, dbpath string) (*query.Queries, error) {
	dir := filepath.Dir(dbpath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "file:"+dbpath+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}

	// Enable foreign key constraints.
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	// Create tables.
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return query.New(db), nil
}
