// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: email.sql

package queries

import (
	"context"
	"database/sql"
)

const countEmails = `-- name: CountEmails :one
SELECT COUNT(*) FROM emails WHERE pid = ?
`

func (q *Queries) CountEmails(ctx context.Context, pid int64) (int64, error) {
	row := q.queryRow(ctx, q.countEmailsStmt, countEmails, pid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createEmail = `-- name: CreateEmail :execresult
INSERT INTO emails (address, pid, verified) VALUES (?, ?, false)
`

type CreateEmailParams struct {
	Address string
	Pid     int64
}

func (q *Queries) CreateEmail(ctx context.Context, arg CreateEmailParams) (sql.Result, error) {
	return q.exec(ctx, q.createEmailStmt, createEmail, arg.Address, arg.Pid)
}

const deleteEmail = `-- name: DeleteEmail :execresult
DELETE FROM emails WHERE id = ?
`

func (q *Queries) DeleteEmail(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteEmailStmt, deleteEmail, id)
}

const getEmail = `-- name: GetEmail :one
SELECT address, created_at, updated_at, verified, pid, id FROM emails WHERE id = ?
`

func (q *Queries) GetEmail(ctx context.Context, id int64) (Email, error) {
	row := q.queryRow(ctx, q.getEmailStmt, getEmail, id)
	var i Email
	err := row.Scan(
		&i.Address,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Verified,
		&i.Pid,
		&i.ID,
	)
	return i, err
}

const getVerifiedEmailByAddress = `-- name: GetVerifiedEmailByAddress :one
SELECT id FROM emails WHERE address = ? AND verified = true
`

func (q *Queries) GetVerifiedEmailByAddress(ctx context.Context, address string) (int64, error) {
	row := q.queryRow(ctx, q.getVerifiedEmailByAddressStmt, getVerifiedEmailByAddress, address)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const listEmails = `-- name: ListEmails :many
SELECT address, created_at, updated_at, verified, pid, id FROM emails WHERE pid = ?
`

func (q *Queries) ListEmails(ctx context.Context, pid int64) ([]Email, error) {
	rows, err := q.query(ctx, q.listEmailsStmt, listEmails, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Email
	for rows.Next() {
		var i Email
		if err := rows.Scan(
			&i.Address,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Verified,
			&i.Pid,
			&i.ID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markEmailVerified = `-- name: MarkEmailVerified :execresult
UPDATE emails SET verified = true WHERE id = ?
`

func (q *Queries) MarkEmailVerified(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.markEmailVerifiedStmt, markEmailVerified, id)
}