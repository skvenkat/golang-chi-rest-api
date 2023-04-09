package repo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AddrBookRepo struct {
	db                               *sqlx.DB
	selectAllContactsWithPhonesStmt  *sqlx.NamedStmt
	insertContactStmt                *sqlx.NamedStmt
	insertPhoneStmt                  *sqlx.NamedStmt
	selectContactsWithPhonesByIdStmt *sqlx.NamedStmt
	deletePhonesByContactIdStmt      *sqlx.NamedStmt
	updateContactByIdStmt            *sqlx.NamedStmt
	deleteContactByIdStmt            *sqlx.NamedStmt
}

func NewAddrBookRepo(db *sqlx.DB) *AddrBookRepo {
	return &AddrBookRepo{
		db:                               db,
		selectAllContactsWithPhonesStmt:  MustPrepareNamed(db, selectAllContactsWithPhonesSql),
		insertContactStmt:                MustPrepareNamed(db, insertContactSql),
		insertPhoneStmt:                  MustPrepareNamed(db, insertPhoneSql),
		selectContactsWithPhonesByIdStmt: MustPrepareNamed(db, selectContactsWithPhonesByIdSql),
		deletePhonesByContactIdStmt:      MustPrepareNamed(db, deletePhonesByContactIdSql),
		updateContactByIdStmt:            MustPrepareNamed(db, updateContactByIdSql),
		deleteContactByIdStmt:            MustPrepareNamed(db, deleteContactByIdSql),
	}
}

// ContactWithPhonesEntity is a result of JOIN
type ContactWithPhonesEntity struct {
	ID        int64
	FirstName string
	LastName  string
	Phones    []*PhoneEntity
}

type PhoneEntity struct {
	PhoneType   string `db:"type"`
	PhoneNumber string `db:"phone_number"`
}

type contactWithPhoneRow struct {
	ID          int64   `db:"id"`
	FirstName   string  `db:"first_name"`
	LastName    string  `db:"last_name"`
	PhoneType   *string `db:"phone_type"`
	PhoneNumber *string `db:"phone_number"`
}

func (p *contactWithPhoneRow) buildPhoneEntity() *PhoneEntity {
	if p == nil || p.PhoneType == nil {
		return nil
	}
	return &PhoneEntity{
		PhoneType:   *p.PhoneType,
		PhoneNumber: *p.PhoneNumber,
	}
}

func (c *contactWithPhoneRow) toContactWithPhonesEntity() *ContactWithPhonesEntity {
	var phones []*PhoneEntity
	phone := c.buildPhoneEntity()
	if phone != nil {
		phones = []*PhoneEntity{phone}
	} else {
		phones = []*PhoneEntity{}
	}
	entity := c.toContactEntity()
	entity.Phones = phones
	return entity
}

func (c *contactWithPhoneRow) toContactEntity() *ContactWithPhonesEntity {
	return &ContactWithPhonesEntity{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
}

func (r *AddrBookRepo) AddContact(ctx context.Context, c *ContactWithPhonesEntity) (*ContactWithPhonesEntity, error) {
	tx := r.db.MustBeginTx(ctx, nil)
	defer tx.Rollback()

	var err error
	newc := *c
	newc.ID, err = ExecNamedStmtReturningLastInsertId(ctx, tx.NamedStmtContext(ctx, r.insertContactStmt), map[string]any{
		"firstName": c.FirstName,
		"lastName":  c.LastName,
	})
	if err != nil {
		err = fmt.Errorf("error inserting contact into database: %w", err)
		zap.S().Errorln(err)
		return nil, err
	}
	// insert phones for contact
	insertPhoneStmt := tx.NamedStmtContext(ctx, r.insertPhoneStmt)
	for _, ph := range c.Phones {
		if err = r.insertPhone(ctx, insertPhoneStmt, newc.ID, ph); err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("error committing transaction: %w", err)
		zap.S().Errorln(err)
		return nil, err
	}
	return &newc, nil
}

func (r *AddrBookRepo) insertPhone(ctx context.Context, insertPhoneStmt *sqlx.NamedStmt, contactId int64, ph *PhoneEntity) error {
	_, err := insertPhoneStmt.ExecContext(ctx, map[string]any{
		"type":        ph.PhoneType,
		"phoneNumber": ph.PhoneNumber,
		"contactId":   contactId,
	})
	if err != nil {
		err = fmt.Errorf("error inserting contact phone into database: %w", err)
		zap.S().Errorln(err)
	}
	return err
}

func (r *AddrBookRepo) UpdateContact(ctx context.Context, c *ContactWithPhonesEntity) (found bool, err error) {
	tx := r.db.MustBeginTx(ctx, nil)
	defer tx.Rollback()

	result, err := tx.NamedStmtContext(ctx, r.updateContactByIdStmt).ExecContext(ctx, map[string]any{
		"contactId": c.ID,
		"firstName": c.FirstName,
		"lastName":  c.LastName,
	})
	if err != nil {
		err = fmt.Errorf("error updating contact id=%d in database: %w", c.ID, err)
		zap.S().Errorln(err)
		return false, err
	}
	if MustGetRowsAffected(result) == 0 {
		zap.S().Warnln("no contact record found by id:", c.ID)
		return false, nil
	}

	_, err = tx.NamedStmtContext(ctx, r.deletePhonesByContactIdStmt).ExecContext(ctx, map[string]any{
		"contactId": c.ID,
	})
	if err != nil {
		err = fmt.Errorf("error updating contact phone in database: %w", err)
		zap.S().Errorln(err)
		return false, err
	}

	// insert phones for contact
	insertPhoneStmt := tx.NamedStmtContext(ctx, r.insertPhoneStmt)
	for _, ph := range c.Phones {
		if err = r.insertPhone(ctx, insertPhoneStmt, c.ID, ph); err != nil {
			return false, err
		}
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("error committing transaction: %w", err)
		zap.S().Errorln(err)
	}
	return true, err
}

func (r *AddrBookRepo) SelectContactByID(ctx context.Context, ID int64) (*ContactWithPhonesEntity, error) {
	var rows []*contactWithPhoneRow
	err := r.selectContactsWithPhonesByIdStmt.SelectContext(ctx, &rows, map[string]any{
		"id": ID,
	})
	if err != nil {
		zap.S().Errorln("Error selecting contact by id=%d in database:", ID, err)
		return nil, err
	}
	if len(rows) == 0 {
		zap.S().Warnf("contact id=%d not found", ID)
		return nil, nil
	}
	row := rows[0]
	entity := row.toContactEntity()
	entity.Phones = make([]*PhoneEntity, len(rows))
	for i, row := range rows {
		entity.Phones[i] = row.buildPhoneEntity()
	}
	return entity, nil
}

func (r *AddrBookRepo) SelectAllContacts(ctx context.Context) ([]*ContactWithPhonesEntity, error) {
	var rows []*contactWithPhoneRow
	err := r.selectAllContactsWithPhonesStmt.SelectContext(ctx, &rows, map[string]any{})
	if err != nil {
		zap.S().Errorln("Error selecting all contacts in database:", err)
		return nil, err
	}
	// The response will be something like:
	// <id> <first_name> <last_name> <"phone.type"> <"phone.phone_number">
	//  2    Toly         Pochkin     home           503-999-9999
	//  3    Julia        Pod         home           333-111-1111
	//  2    Toly         Pochkin     mobile         503-555-7777
	//  3    Julia        Pod         mobile         333-555-2222
	// Our goal is to get rid of duplicate contacts while merging phones from duplicate contacts into a single distinct contact
	entities := MergeJoinedRows[*contactWithPhoneRow, *ContactWithPhonesEntity](
		rows,
		/*getId*/ func(c *contactWithPhoneRow) int64 {
			return c.ID
		},
		/*createMerged*/ func(row *contactWithPhoneRow) *ContactWithPhonesEntity {
			return row.toContactWithPhonesEntity()
		},
		/*updateMerged*/ func(row *contactWithPhoneRow, existingEntity *ContactWithPhonesEntity) *ContactWithPhonesEntity {
			existingEntity.Phones = append(existingEntity.Phones, row.buildPhoneEntity())
			return existingEntity
		},
	)
	return entities, nil
}

func (r *AddrBookRepo) DeleteContact(ctx context.Context, id int64) (found bool, err error) {
	tx := r.db.MustBeginTx(ctx, nil)
	defer tx.Rollback()

	_, err = tx.NamedStmtContext(ctx, r.deletePhonesByContactIdStmt).ExecContext(ctx, map[string]any{
		"contactId": id,
	})
	if err != nil {
		err = fmt.Errorf("error deleting phone contacts by contact id=%d: %w", id, err)
		zap.S().Errorln(err)
		return false, err
	}

	result, err := tx.NamedStmtContext(ctx, r.deleteContactByIdStmt).ExecContext(ctx, map[string]any{
		"id": id,
	})
	found = MustGetRowsAffected(result) > 0

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("error committing transaction: %w", err)
		zap.S().Errorln(err)
	}
	return
}
