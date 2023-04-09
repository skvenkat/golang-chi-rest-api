package outport

import "github.com/jmoiron/sqlx"

type Persistence interface {
	DB() *sqlx.DB
	Close()
}
