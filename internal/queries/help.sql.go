// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: help.sql

package queries

import (
	"context"
)

const getHelp = `-- name: GetHelp :one
SELECT created_at, updated_at, html, raw, sub, title, slug, pid FROM help WHERE slug = ?
`

func (q *Queries) GetHelp(ctx context.Context, slug string) (Help, error) {
	row := q.queryRow(ctx, q.getHelpStmt, getHelp, slug)
	var i Help
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HTML,
		&i.Raw,
		&i.Sub,
		&i.Title,
		&i.Slug,
		&i.PID,
	)
	return i, err
}

const getHelpRelated = `-- name: GetHelpRelated :many
SELECT related_sub, related_title, related, slug FROM help_related WHERE slug = ?
`

func (q *Queries) GetHelpRelated(ctx context.Context, slug string) ([]HelpRelated, error) {
	rows, err := q.query(ctx, q.getHelpRelatedStmt, getHelpRelated, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []HelpRelated
	for rows.Next() {
		var i HelpRelated
		if err := rows.Scan(
			&i.RelatedSub,
			&i.RelatedTitle,
			&i.Related,
			&i.Slug,
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

const listHelpSlugs = `-- name: ListHelpSlugs :many
SELECT slug FROM help
`

func (q *Queries) ListHelpSlugs(ctx context.Context) ([]string, error) {
	rows, err := q.query(ctx, q.listHelpSlugsStmt, listHelpSlugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		items = append(items, slug)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listHelpTitleAndSub = `-- name: ListHelpTitleAndSub :many
SELECT slug, title, sub FROM help ORDER BY slug
`

type ListHelpTitleAndSubRow struct {
	Slug  string
	Title string
	Sub   string
}

func (q *Queries) ListHelpTitleAndSub(ctx context.Context) ([]ListHelpTitleAndSubRow, error) {
	rows, err := q.query(ctx, q.listHelpTitleAndSubStmt, listHelpTitleAndSub)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListHelpTitleAndSubRow
	for rows.Next() {
		var i ListHelpTitleAndSubRow
		if err := rows.Scan(&i.Slug, &i.Title, &i.Sub); err != nil {
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