package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func ExecNamedStmtReturningLastInsertId(ctx context.Context, stmt *sqlx.NamedStmt, arg any) (int64, error) {
	result, err := stmt.ExecContext(ctx, arg)
	if err == nil {
		return result.LastInsertId()
	}
	return 0, err
}
