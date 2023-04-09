package persist

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/outport"
	"go.uber.org/zap"
)

const createTablesSql =
/*language=sqlite*/ `
CREATE TABLE IF NOT EXISTS contacts(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS phones(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    contact_id BIGINT NOT NULL REFERENCES contacts(id)
)
`

type dbAdapter struct {
	db *sqlx.DB
}

// NewPersistence connects to SQLite database and returns Persistence interface that wraps database reference
func NewPersistence(cfg *app.Config) outport.Persistence {
	dbcfg := cfg.Database
	connStr := fmt.Sprintf("%s", dbcfg.Filename)
	zap.S().Infoln("establishing connection to SQLite database...")
	db, err := sqlx.ConnectContext(context.Background(), "sqlite3", connStr)
	if err != nil {
		zap.S().Fatalf("error connecting to database (file=%s): %s\n", dbcfg.Filename, err)
	}
	zap.S().Infoln("connection to database was successfully established, performing initialization...")

	tx := db.MustBegin()
	defer tx.Rollback()
	tx.MustExec(createTablesSql)
	err = tx.Commit()
	if err != nil {
		zap.S().Fatalln("failed to commit transaction while creating database:", err)
	}
	zap.S().Infoln("db initialization was successfully performed")

	return &dbAdapter{db: db}
}

func (d dbAdapter) DB() *sqlx.DB {
	return d.db
}

func (d dbAdapter) Close() {
	err := d.DB().Close()
	if err != nil {
		fmt.Println("failed to close sqlite database:", err)
	}
}
