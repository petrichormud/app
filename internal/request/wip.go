package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
)

// TODO: Find a way to get each content extractor into the definition of the request type
func GetContent(qtx *queries.Queries, req *queries.Request) (map[string]string, error) {
	var b []byte
	m := map[string]string{}

	switch req.Type {
	case TypeCharacterApplication:
		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), req.ID)
		if err != nil {
			return m, err
		}

		b, err = json.Marshal(app)
		if err != nil {
			return m, err
		}
	default:
		return m, errors.New("invalid type")
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]string{}, err
	}

	return m, nil
}

func GetNextIncompleteField(t string, content map[string]string) (string, bool) {
	fields := FieldNamesByType[t]
	for i, field := range fields {
		value, ok := content[field]
		if !ok {
			continue
		}
		if len(value) == 0 {
			return field, i == len(fields)-1
		}
	}
	return "", false
}

// TODO: Make this FieldTag?
type UpdateFieldParams struct {
	Request *queries.Request
	Field   string
	Value   string
	PID     int64
}

// TODO: Turn this into a map of updaters by field - can create an interface for the Updater
func UpdateField(qtx *queries.Queries, p UpdateFieldParams) error {
	switch p.Field {
	case FieldName:
		if !IsNameValid(p.Value) {
			return ErrInvalidInput
		}

		if err := qtx.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			RID:  p.Request.ID,
			Name: p.Value,
		}); err != nil {
			return err
		}

	case FieldGender:
		if !IsGenderValid(p.Value) {
			return ErrInvalidInput
		}

		if err := qtx.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			RID:    p.Request.ID,
			Gender: p.Value,
		}); err != nil {
			return err
		}
	case FieldShortDescription:
		if !IsShortDescriptionValid(p.Value) {
			return ErrInvalidInput
		}

		if err := qtx.UpdateCharacterApplicationContentShortDescription(context.Background(), queries.UpdateCharacterApplicationContentShortDescriptionParams{
			RID:              p.Request.ID,
			ShortDescription: p.Value,
		}); err != nil {
			return err
		}
	case FieldDescription:
		if !IsDescriptionValid(p.Value) {
			return ErrInvalidInput
		}

		if err := qtx.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			RID:         p.Request.ID,
			Description: p.Value,
		}); err != nil {
			return err
		}
	case FieldBackstory:
		if !IsBackstoryValid(p.Value) {
			return ErrInvalidInput
		}

		if err := qtx.UpdateCharacterApplicationContentBackstory(context.Background(), queries.UpdateCharacterApplicationContentBackstoryParams{
			RID:       p.Request.ID,
			Backstory: p.Value,
		}); err != nil {
			return err
		}
	default:
		return ErrMalformedUpdateInput
	}

	if p.Request.Type == TypeCharacterApplication {
		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), p.Request.ID)
		if err != nil {
			return err
		}

		ready := IsCharacterApplicationValid(&app)

		if ready && p.Request.Status == StatusIncomplete {
			if err := qtx.CreateHistoryForRequestStatusChange(context.Background(), queries.CreateHistoryForRequestStatusChangeParams{
				RID: p.Request.ID,
				PID: p.PID,
			}); err != nil {
				return err
			}

			if err := qtx.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
				ID:     p.Request.ID,
				Status: StatusReady,
			}); err != nil {
				return err
			}
		} else if !ready && p.Request.Status == StatusReady {
			if err := qtx.CreateHistoryForRequestStatusChange(context.Background(), queries.CreateHistoryForRequestStatusChangeParams{
				RID: p.Request.ID,
				PID: p.PID,
			}); err != nil {
				return err
			}

			if err := qtx.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
				ID:     p.Request.ID,
				Status: StatusIncomplete,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: Get this built into the Definition
func GetSummaryTitle(t string, content map[string]string) string {
	if t == TypeCharacterApplication {
		var sb strings.Builder
		titleName := constants.DefaultName
		if len(content[FieldName]) > 0 {
			titleName = content[FieldName]
		}
		fmt.Fprintf(&sb, "Character Application (%s)", titleName)
		return sb.String()
	}

	return "Request"
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
}

type GetSummaryFieldsParams struct {
	Request *queries.Request
	Content map[string]string
	PID     int64
}

// TODO: Get this built into the Definition
func GetSummaryFields(p GetSummaryFieldsParams) []SummaryField {
	if p.Request.Type == TypeCharacterApplication {
		var basePathSB strings.Builder
		fmt.Fprintf(&basePathSB, "/requests/%d", p.Request.ID)
		basePath := basePathSB.String()

		var namePathSB strings.Builder
		fmt.Fprintf(&namePathSB, "%s/%s", basePath, FieldName)

		var genderPathSB strings.Builder
		fmt.Fprintf(&genderPathSB, "%s/%s", basePath, FieldGender)

		var shortDescriptionPathSB strings.Builder
		fmt.Fprintf(&shortDescriptionPathSB, "%s/%s", basePath, FieldShortDescription)

		var descriptionPathSB strings.Builder
		fmt.Fprintf(&descriptionPathSB, "%s/%s", basePath, FieldDescription)

		var backstoryPathSB strings.Builder
		fmt.Fprintf(&backstoryPathSB, "%s/%s", basePath, FieldBackstory)

		allowEdit := p.Request.PID == p.PID

		return []SummaryField{
			{
				Label:     "Name",
				Content:   p.Content[FieldName],
				AllowEdit: allowEdit,
				Path:      namePathSB.String(),
			},
			{
				Label:     "Gender",
				Content:   p.Content[FieldGender],
				AllowEdit: allowEdit,
				Path:      genderPathSB.String(),
			},
			{
				Label:     "Short Description",
				Content:   p.Content[FieldShortDescription],
				AllowEdit: allowEdit,
				Path:      shortDescriptionPathSB.String(),
			},
			{
				Label:     "Description",
				Content:   p.Content[FieldDescription],
				AllowEdit: allowEdit,
				Path:      descriptionPathSB.String(),
			},
			{
				Label:     "Backstory",
				Content:   p.Content[FieldBackstory],
				AllowEdit: allowEdit,
				Path:      backstoryPathSB.String(),
			},
		}
	}

	return []SummaryField{}
}
