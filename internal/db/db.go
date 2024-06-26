package db

import (
	"embed"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func New(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return db, err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectSQLite3)); err != nil {
		return db, err
	}

	rawDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if err := goose.Up(rawDB, "migrations"); err != nil {
		return db, err
	}

	return db, err
}
