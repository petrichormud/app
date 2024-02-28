package request

import (
	"encoding/json"
	"errors"

	"petrichormud.com/app/internal/query"
)

const errNoDefinition string = "no definition with type"

var ErrNoDefinition error = errors.New(errNoDefinition)

type Definition interface {
	Type() string
	Dialogs() Dialogs
	Fields() Fields
	IsFieldNameValid(f string) bool
	Content(q *query.Queries, rid int64) (content, error)
	IsContentValid(c content) bool
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	SummaryTitle(content map[string]string) string
	SummaryFields(p GetSummaryFieldsParams) []SummaryField
}

type DefaultDefinition struct{}

func (d *DefaultDefinition) Fields() Fields {
	return Fields{}
}

func (d *DefaultDefinition) Content(q *query.Queries, rid int64) (content, error) {
	return content{}, ErrNoDefinition
}

func (d *DefaultDefinition) IsContentValid(c content) bool {
	return false
}

func (d *DefaultDefinition) UpdateField(q *query.Queries, p UpdateFieldParams) error {
	fields := d.Fields()
	if err := fields.Update(q, p); err != nil {
		return err
	}

	c, err := d.Content(q, p.Request.ID)
	if err != nil {
		return err
	}
	ready := d.IsContentValid(c)

	if err := UpdateReadyStatus(q, UpdateReadyStatusParams{
		Status: p.Request.Status,
		PID:    p.PID,
		RID:    p.Request.ID,
		Ready:  ready,
	}); err != nil {
		return err
	}

	return nil
}

type UpdateFieldParams struct {
	Request   *query.Request
	FieldName string
	Value     string
	PID       int64
}

func UpdateField(q *query.Queries, p UpdateFieldParams) error {
	if !IsTypeValid(p.Request.Type) {
		return ErrInvalidType
	}
	definition, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return ErrNoDefinition
	}
	if err := definition.UpdateField(q, p); err != nil {
		return err
	}
	return nil
}

func View(t, f string) string {
	definition, ok := Definitions.Get(t)
	if !ok {
		return ""
	}
	fields := definition.Fields().Map
	field := fields[f]
	return field.View
}

func ContentBytes(content any) ([]byte, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// TODO: Let this return the fully-qualified type
func Content(q *query.Queries, req *query.Request) (map[string]string, error) {
	m := map[string]string{}

	if !IsTypeValid(req.Type) {
		return m, ErrInvalidType
	}

	definition, ok := Definitions.Get(req.Type)
	if !ok {
		return m, ErrNoDefinition
	}

	content, err := definition.Content(q, req.ID)
	if err != nil {
		return m, err
	}

	return content.Inner, nil
}

// TODO: Clean this up based on the Fields or new Content API
func NextIncompleteField(t string, content map[string]string) (string, bool) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", false
	}
	fields := definition.Fields()
	for i, field := range fields.List {
		value, ok := content[field.Name]
		if !ok {
			continue
		}
		if len(value) == 0 {
			return field.Name, i == len(fields.List)-1
		}
	}
	return "", false
}

// TODO: Possible error output
func GetFieldLabelAndDescription(t, f string) (string, string) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", ""
	}
	fields := definition.Fields().Map
	field := fields[f]
	return field.Label, field.Description
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
}

type GetSummaryFieldsParams struct {
	Request *query.Request
	Content map[string]string
	PID     int64
}

func SummaryFields(p GetSummaryFieldsParams) []SummaryField {
	definition, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return []SummaryField{}
	}
	return definition.SummaryFields(p)
}

func SummaryTitle(t string, content map[string]string) string {
	if !IsTypeValid(t) {
		return "Request"
	}
	definition, ok := Definitions.Get(t)
	if !ok {
		return ""
	}
	return definition.SummaryTitle(content)
}

func IsFieldNameValid(t, name string) bool {
	definition, ok := Definitions.Get(t)
	if !ok {
		return false
	}
	return definition.IsFieldNameValid(name)
}
