// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package odb

import (
	"database/sql"
	"encoding/json"
	"time"
)

type AtlasSchemaRevision struct {
	Version         string
	Description     string
	Type            uint64
	Applied         int64
	Total           int64
	ExecutedAt      time.Time
	ExecutionTime   int64
	Error           sql.NullString
	ErrorStmt       sql.NullString
	Hash            string
	PartialHashes   json.RawMessage
	OperatorVersion string
}

type Feature struct {
	ID   int64
	Flag string
}

type Player struct {
	ID       int64
	Username string
	PwHash   string
}

type Request struct {
	ID int64
}
