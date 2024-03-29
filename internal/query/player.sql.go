// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: player.sql

package query

import (
	"context"
	"database/sql"
)

const createPlayer = `-- name: CreatePlayer :execresult
INSERT INTO players (username, pw_hash) VALUES (?, ?)
`

type CreatePlayerParams struct {
	Username string
	PwHash   string
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (sql.Result, error) {
	return q.exec(ctx, q.createPlayerStmt, createPlayer, arg.Username, arg.PwHash)
}

const createPlayerPermission = `-- name: CreatePlayerPermission :execresult
INSERT INTO player_permissions (name, pid, ipid) VALUES (?, ?, ?)
`

type CreatePlayerPermissionParams struct {
	Name string
	PID  int64
	IPID int64
}

func (q *Queries) CreatePlayerPermission(ctx context.Context, arg CreatePlayerPermissionParams) (sql.Result, error) {
	return q.exec(ctx, q.createPlayerPermissionStmt, createPlayerPermission, arg.Name, arg.PID, arg.IPID)
}

const createPlayerPermissionIssuedChangeHistory = `-- name: CreatePlayerPermissionIssuedChangeHistory :exec
INSERT INTO player_permission_change_history (name, pid, ipid) VALUES (?, ?, ?)
`

type CreatePlayerPermissionIssuedChangeHistoryParams struct {
	Name string
	PID  int64
	IPID int64
}

func (q *Queries) CreatePlayerPermissionIssuedChangeHistory(ctx context.Context, arg CreatePlayerPermissionIssuedChangeHistoryParams) error {
	_, err := q.exec(ctx, q.createPlayerPermissionIssuedChangeHistoryStmt, createPlayerPermissionIssuedChangeHistory, arg.Name, arg.PID, arg.IPID)
	return err
}

const createPlayerPermissionRevokedChangeHistory = `-- name: CreatePlayerPermissionRevokedChangeHistory :exec
INSERT INTO player_permission_change_history (name, pid, ipid, revoked) VALUES (?, ?, ?, true)
`

type CreatePlayerPermissionRevokedChangeHistoryParams struct {
	Name string
	PID  int64
	IPID int64
}

func (q *Queries) CreatePlayerPermissionRevokedChangeHistory(ctx context.Context, arg CreatePlayerPermissionRevokedChangeHistoryParams) error {
	_, err := q.exec(ctx, q.createPlayerPermissionRevokedChangeHistoryStmt, createPlayerPermissionRevokedChangeHistory, arg.Name, arg.PID, arg.IPID)
	return err
}

const createPlayerSettings = `-- name: CreatePlayerSettings :exec
INSERT INTO player_settings (theme, pid) VALUES (?, ?)
`

type CreatePlayerSettingsParams struct {
	Theme string
	PID   int64
}

func (q *Queries) CreatePlayerSettings(ctx context.Context, arg CreatePlayerSettingsParams) error {
	_, err := q.exec(ctx, q.createPlayerSettingsStmt, createPlayerSettings, arg.Theme, arg.PID)
	return err
}

const deletePlayerPermission = `-- name: DeletePlayerPermission :exec
DELETE FROM player_permissions WHERE name = ? AND pid = ?
`

type DeletePlayerPermissionParams struct {
	Name string
	PID  int64
}

func (q *Queries) DeletePlayerPermission(ctx context.Context, arg DeletePlayerPermissionParams) error {
	_, err := q.exec(ctx, q.deletePlayerPermissionStmt, deletePlayerPermission, arg.Name, arg.PID)
	return err
}

const getPlayer = `-- name: GetPlayer :one
SELECT created_at, updated_at, pw_hash, username, id FROM players WHERE id = ?
`

func (q *Queries) GetPlayer(ctx context.Context, id int64) (Player, error) {
	row := q.queryRow(ctx, q.getPlayerStmt, getPlayer, id)
	var i Player
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PwHash,
		&i.Username,
		&i.ID,
	)
	return i, err
}

const getPlayerByUsername = `-- name: GetPlayerByUsername :one
SELECT created_at, updated_at, pw_hash, username, id FROM players WHERE username = ?
`

func (q *Queries) GetPlayerByUsername(ctx context.Context, username string) (Player, error) {
	row := q.queryRow(ctx, q.getPlayerByUsernameStmt, getPlayerByUsername, username)
	var i Player
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PwHash,
		&i.Username,
		&i.ID,
	)
	return i, err
}

const getPlayerSettings = `-- name: GetPlayerSettings :one
SELECT created_at, updated_at, theme, pid, id FROM player_settings WHERE pid = ?
`

func (q *Queries) GetPlayerSettings(ctx context.Context, pid int64) (PlayerSetting, error) {
	row := q.queryRow(ctx, q.getPlayerSettingsStmt, getPlayerSettings, pid)
	var i PlayerSetting
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Theme,
		&i.PID,
		&i.ID,
	)
	return i, err
}

const getPlayerUsername = `-- name: GetPlayerUsername :one
SELECT (username) FROM players WHERE id = ?
`

func (q *Queries) GetPlayerUsername(ctx context.Context, id int64) (string, error) {
	row := q.queryRow(ctx, q.getPlayerUsernameStmt, getPlayerUsername, id)
	var username string
	err := row.Scan(&username)
	return username, err
}

const getPlayerUsernameById = `-- name: GetPlayerUsernameById :one
SELECT (username) FROM players WHERE id = ?
`

func (q *Queries) GetPlayerUsernameById(ctx context.Context, id int64) (string, error) {
	row := q.queryRow(ctx, q.getPlayerUsernameByIdStmt, getPlayerUsernameById, id)
	var username string
	err := row.Scan(&username)
	return username, err
}

const listPlayerPermissions = `-- name: ListPlayerPermissions :many
SELECT created_at, name, ipid, pid, id FROM player_permissions WHERE pid = ?
`

func (q *Queries) ListPlayerPermissions(ctx context.Context, pid int64) ([]PlayerPermission, error) {
	rows, err := q.query(ctx, q.listPlayerPermissionsStmt, listPlayerPermissions, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PlayerPermission
	for rows.Next() {
		var i PlayerPermission
		if err := rows.Scan(
			&i.CreatedAt,
			&i.Name,
			&i.IPID,
			&i.PID,
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

const searchPlayersByUsername = `-- name: SearchPlayersByUsername :many
SELECT created_at, updated_at, pw_hash, username, id FROM players WHERE username LIKE ?
`

func (q *Queries) SearchPlayersByUsername(ctx context.Context, username string) ([]Player, error) {
	rows, err := q.query(ctx, q.searchPlayersByUsernameStmt, searchPlayersByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Player
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PwHash,
			&i.Username,
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

const updatePlayerPassword = `-- name: UpdatePlayerPassword :execresult
UPDATE players SET pw_hash = ? WHERE id = ?
`

type UpdatePlayerPasswordParams struct {
	PwHash string
	ID     int64
}

func (q *Queries) UpdatePlayerPassword(ctx context.Context, arg UpdatePlayerPasswordParams) (sql.Result, error) {
	return q.exec(ctx, q.updatePlayerPasswordStmt, updatePlayerPassword, arg.PwHash, arg.ID)
}

const updatePlayerSettingsTheme = `-- name: UpdatePlayerSettingsTheme :exec
UPDATE player_settings SET theme = ? WHERE pid = ?
`

type UpdatePlayerSettingsThemeParams struct {
	Theme string
	PID   int64
}

func (q *Queries) UpdatePlayerSettingsTheme(ctx context.Context, arg UpdatePlayerSettingsThemeParams) error {
	_, err := q.exec(ctx, q.updatePlayerSettingsThemeStmt, updatePlayerSettingsTheme, arg.Theme, arg.PID)
	return err
}
