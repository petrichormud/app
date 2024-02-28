// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: help.sql

package query

import (
	"context"
)

const getHelp = `-- name: GetHelp :one
SELECT created_at, updated_at, html, raw, sub, title, category, slug, pid FROM help WHERE slug = ?
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
		&i.Category,
		&i.Slug,
		&i.PID,
	)
	return i, err
}

const getHelpRelated = `-- name: GetHelpRelated :many
SELECT related_sub, related_title, related_slug, slug FROM help_related WHERE slug = ?
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
			&i.RelatedSlug,
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

const getTagsForHelpFile = `-- name: GetTagsForHelpFile :many
SELECT tag FROM help_tags WHERE slug = ?
`

func (q *Queries) GetTagsForHelpFile(ctx context.Context, slug string) ([]string, error) {
	rows, err := q.query(ctx, q.getTagsForHelpFileStmt, getTagsForHelpFile, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		items = append(items, tag)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listHelpHeaders = `-- name: ListHelpHeaders :many
SELECT slug, title, sub, category FROM help ORDER BY slug
`

type ListHelpHeadersRow struct {
	Slug     string
	Title    string
	Sub      string
	Category string
}

func (q *Queries) ListHelpHeaders(ctx context.Context) ([]ListHelpHeadersRow, error) {
	rows, err := q.query(ctx, q.listHelpHeadersStmt, listHelpHeaders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListHelpHeadersRow
	for rows.Next() {
		var i ListHelpHeadersRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Sub,
			&i.Category,
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

const searchHelpByCategory = `-- name: SearchHelpByCategory :many
SELECT slug, title, sub, category FROM help WHERE category LIKE ?
`

type SearchHelpByCategoryRow struct {
	Slug     string
	Title    string
	Sub      string
	Category string
}

func (q *Queries) SearchHelpByCategory(ctx context.Context, category string) ([]SearchHelpByCategoryRow, error) {
	rows, err := q.query(ctx, q.searchHelpByCategoryStmt, searchHelpByCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchHelpByCategoryRow
	for rows.Next() {
		var i SearchHelpByCategoryRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Sub,
			&i.Category,
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

const searchHelpByContent = `-- name: SearchHelpByContent :many
SELECT slug, title, sub, category FROM help WHERE raw LIKE ? OR sub LIKE ?
`

type SearchHelpByContentParams struct {
	Raw string
	Sub string
}

type SearchHelpByContentRow struct {
	Slug     string
	Title    string
	Sub      string
	Category string
}

func (q *Queries) SearchHelpByContent(ctx context.Context, arg SearchHelpByContentParams) ([]SearchHelpByContentRow, error) {
	rows, err := q.query(ctx, q.searchHelpByContentStmt, searchHelpByContent, arg.Raw, arg.Sub)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchHelpByContentRow
	for rows.Next() {
		var i SearchHelpByContentRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Sub,
			&i.Category,
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

const searchHelpByTags = `-- name: SearchHelpByTags :many
SELECT slug, title, sub, category FROM help WHERE slug IN (SELECT slug FROM help_tags WHERE tag LIKE ?)
`

type SearchHelpByTagsRow struct {
	Slug     string
	Title    string
	Sub      string
	Category string
}

func (q *Queries) SearchHelpByTags(ctx context.Context, tag string) ([]SearchHelpByTagsRow, error) {
	rows, err := q.query(ctx, q.searchHelpByTagsStmt, searchHelpByTags, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchHelpByTagsRow
	for rows.Next() {
		var i SearchHelpByTagsRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Sub,
			&i.Category,
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

const searchHelpByTitle = `-- name: SearchHelpByTitle :many
SELECT slug, title, sub, category FROM help WHERE slug LIKE ? OR title LIKE ?
`

type SearchHelpByTitleParams struct {
	Slug  string
	Title string
}

type SearchHelpByTitleRow struct {
	Slug     string
	Title    string
	Sub      string
	Category string
}

func (q *Queries) SearchHelpByTitle(ctx context.Context, arg SearchHelpByTitleParams) ([]SearchHelpByTitleRow, error) {
	rows, err := q.query(ctx, q.searchHelpByTitleStmt, searchHelpByTitle, arg.Slug, arg.Title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchHelpByTitleRow
	for rows.Next() {
		var i SearchHelpByTitleRow
		if err := rows.Scan(
			&i.Slug,
			&i.Title,
			&i.Sub,
			&i.Category,
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

const searchTags = `-- name: SearchTags :many
SELECT slug, tag FROM help_tags WHERE tag LIKE ?
`

type SearchTagsRow struct {
	Slug string
	Tag  string
}

func (q *Queries) SearchTags(ctx context.Context, tag string) ([]SearchTagsRow, error) {
	rows, err := q.query(ctx, q.searchTagsStmt, searchTags, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchTagsRow
	for rows.Next() {
		var i SearchTagsRow
		if err := rows.Scan(&i.Slug, &i.Tag); err != nil {
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