package repo

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// MustPrepareNamed creates prepared statement or panics if it cannot be created.
// Usually you want to create named statement when application starts.
func MustPrepareNamed(db *sqlx.DB, query string) *sqlx.NamedStmt {
	if stmt, err := db.PrepareNamed(query); err == nil {
		return stmt
	} else {
		panic(err)
	}
}

// MergeJoinedRows allows deduplicate joining rows into rows with merged values
func MergeJoinedRows[TJoined any, TMerged any, TId comparable](
	rows []TJoined,
	getId func(row TJoined) TId,
	createMerged func(row TJoined) TMerged,
	updateMerged func(row TJoined, existingEntity TMerged) TMerged,
) []TMerged {
	entities := make([]TMerged, 0, len(rows))
	entityIndexByRowID := make(map[TId]int)
	outputIdx := 0
	for _, row := range rows {
		id := getId(row)
		entityIndex, ok := entityIndexByRowID[id]
		if ok {
			existingRow := entities[entityIndex]
			entities[entityIndex] = updateMerged(row, existingRow)
		} else {
			entity := createMerged(row)
			entities = append(entities, entity)
			entityIndexByRowID[id] = outputIdx
			outputIdx += 1
		}
	}
	return entities
}

func MustGetRowsAffected(result sql.Result) int64 {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Sprintf("error retrieving number of affected rows: %v", err))
	}
	return rowsAffected
}
