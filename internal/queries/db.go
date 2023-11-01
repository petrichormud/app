// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package queries

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.countPlayerEmailsStmt, err = db.PrepareContext(ctx, countPlayerEmails); err != nil {
		return nil, fmt.Errorf("error preparing query CountPlayerEmails: %w", err)
	}
	if q.createPlayerStmt, err = db.PrepareContext(ctx, createPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayer: %w", err)
	}
	if q.createPlayerEmailStmt, err = db.PrepareContext(ctx, createPlayerEmail); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerEmail: %w", err)
	}
	if q.createPlayerPermissionsStmt, err = db.PrepareContext(ctx, createPlayerPermissions); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerPermissions: %w", err)
	}
	if q.getEmailStmt, err = db.PrepareContext(ctx, getEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetEmail: %w", err)
	}
	if q.getPlayerStmt, err = db.PrepareContext(ctx, getPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayer: %w", err)
	}
	if q.getPlayerByUsernameStmt, err = db.PrepareContext(ctx, getPlayerByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerByUsername: %w", err)
	}
	if q.getPlayerPWHashStmt, err = db.PrepareContext(ctx, getPlayerPWHash); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerPWHash: %w", err)
	}
	if q.getPlayerUsernameStmt, err = db.PrepareContext(ctx, getPlayerUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerUsername: %w", err)
	}
	if q.listPlayerEmailsStmt, err = db.PrepareContext(ctx, listPlayerEmails); err != nil {
		return nil, fmt.Errorf("error preparing query ListPlayerEmails: %w", err)
	}
	if q.listPlayerPermissionsStmt, err = db.PrepareContext(ctx, listPlayerPermissions); err != nil {
		return nil, fmt.Errorf("error preparing query ListPlayerPermissions: %w", err)
	}
	if q.markEmailVerifiedStmt, err = db.PrepareContext(ctx, markEmailVerified); err != nil {
		return nil, fmt.Errorf("error preparing query MarkEmailVerified: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.countPlayerEmailsStmt != nil {
		if cerr := q.countPlayerEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countPlayerEmailsStmt: %w", cerr)
		}
	}
	if q.createPlayerStmt != nil {
		if cerr := q.createPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerStmt: %w", cerr)
		}
	}
	if q.createPlayerEmailStmt != nil {
		if cerr := q.createPlayerEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerEmailStmt: %w", cerr)
		}
	}
	if q.createPlayerPermissionsStmt != nil {
		if cerr := q.createPlayerPermissionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerPermissionsStmt: %w", cerr)
		}
	}
	if q.getEmailStmt != nil {
		if cerr := q.getEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEmailStmt: %w", cerr)
		}
	}
	if q.getPlayerStmt != nil {
		if cerr := q.getPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStmt: %w", cerr)
		}
	}
	if q.getPlayerByUsernameStmt != nil {
		if cerr := q.getPlayerByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerByUsernameStmt: %w", cerr)
		}
	}
	if q.getPlayerPWHashStmt != nil {
		if cerr := q.getPlayerPWHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerPWHashStmt: %w", cerr)
		}
	}
	if q.getPlayerUsernameStmt != nil {
		if cerr := q.getPlayerUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerUsernameStmt: %w", cerr)
		}
	}
	if q.listPlayerEmailsStmt != nil {
		if cerr := q.listPlayerEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listPlayerEmailsStmt: %w", cerr)
		}
	}
	if q.listPlayerPermissionsStmt != nil {
		if cerr := q.listPlayerPermissionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listPlayerPermissionsStmt: %w", cerr)
		}
	}
	if q.markEmailVerifiedStmt != nil {
		if cerr := q.markEmailVerifiedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing markEmailVerifiedStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                          DBTX
	tx                          *sql.Tx
	countPlayerEmailsStmt       *sql.Stmt
	createPlayerStmt            *sql.Stmt
	createPlayerEmailStmt       *sql.Stmt
	createPlayerPermissionsStmt *sql.Stmt
	getEmailStmt                *sql.Stmt
	getPlayerStmt               *sql.Stmt
	getPlayerByUsernameStmt     *sql.Stmt
	getPlayerPWHashStmt         *sql.Stmt
	getPlayerUsernameStmt       *sql.Stmt
	listPlayerEmailsStmt        *sql.Stmt
	listPlayerPermissionsStmt   *sql.Stmt
	markEmailVerifiedStmt       *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                          tx,
		tx:                          tx,
		countPlayerEmailsStmt:       q.countPlayerEmailsStmt,
		createPlayerStmt:            q.createPlayerStmt,
		createPlayerEmailStmt:       q.createPlayerEmailStmt,
		createPlayerPermissionsStmt: q.createPlayerPermissionsStmt,
		getEmailStmt:                q.getEmailStmt,
		getPlayerStmt:               q.getPlayerStmt,
		getPlayerByUsernameStmt:     q.getPlayerByUsernameStmt,
		getPlayerPWHashStmt:         q.getPlayerPWHashStmt,
		getPlayerUsernameStmt:       q.getPlayerUsernameStmt,
		listPlayerEmailsStmt:        q.listPlayerEmailsStmt,
		listPlayerPermissionsStmt:   q.listPlayerPermissionsStmt,
		markEmailVerifiedStmt:       q.markEmailVerifiedStmt,
	}
}
