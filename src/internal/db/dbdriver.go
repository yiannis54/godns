package db

import (
	"database/sql"
	"fmt"
	"time"

	migrate "github.com/rubenv/sql-migrate"
)

type DB struct {
	SQL *sql.DB
}

const (
	maxOpenDBConn = 10
	maxIdleDBConn = 5
	maxDBLifetime = 5 * time.Minute
)

func InitDB(dbDSN string) (*DB, error) {
	sqlDB, err := NewSQLite3(dbDSN)
	if err != nil {
		return nil, err
	}

	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "1",
				Up: []string{
					`CREATE TABLE dns_updates (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						new_ip string,
						created_at timestamp);`,
				},
				Down: []string{"DROP TABLE dns_updates"},
			},
		},
	}

	n, err := migrate.Exec(sqlDB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Applied %d migrations!\n", n)

	sqlDB.SetMaxOpenConns(maxOpenDBConn)
	sqlDB.SetMaxIdleConns(maxIdleDBConn)
	sqlDB.SetConnMaxLifetime(maxDBLifetime)

	return &DB{
		SQL: sqlDB,
	}, nil
}

func NewSQLite3(dbDSN string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbDSN))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
