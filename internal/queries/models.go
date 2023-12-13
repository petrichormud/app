// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package queries

import (
	"database/sql"
)

type CharacterApplicationContent struct {
	CreatedAt        sql.NullTime
	UpdatedAt        sql.NullTime
	Backstory        string
	Description      string
	ShortDescription string
	Name             string
	Gender           string
	RID              int64
	ID               int64
}

type CharacterApplicationContentHistory struct {
	CreatedAt        sql.NullTime
	UpdatedAt        sql.NullTime
	Backstory        string
	Description      string
	ShortDescription string
	Name             string
	Gender           string
	RID              int64
	ID               int64
	VID              int32
}

type Email struct {
	Address   string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Verified  bool
	PID       int64
	ID        int64
}

type Player struct {
	PwHash    string
	Username  string
	Role      string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	ID        int64
}

type Request struct {
	Type      string
	Status    string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	PID       int64
	ID        int64
	VID       int32
	New       bool
}

type RequestComment struct {
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	Text      string
	Field     string
	RID       int64
	PID       int64
	CID       int64
	ID        int64
	VID       int32
	Deleted   sql.NullBool
}

type RequestCommentHistory struct {
	CreatedAt sql.NullTime
	Text      string
	Field     string
	CID       int64
	ID        int64
	VID       int32
}

type RequestStatusChange struct {
	Status    string
	CreatedAt sql.NullTime
	RID       int64
	PID       int64
	ID        int64
	VID       int32
}
