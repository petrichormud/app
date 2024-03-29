package request

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
)

const (
	StatusIncomplete = "Incomplete"
	StatusReady      = "Ready"
	StatusSubmitted  = "Submitted"
	StatusInReview   = "InReview"
	StatusApproved   = "Approved"
	StatusReviewed   = "Reviewed"
	StatusRejected   = "Rejected"
	StatusArchived   = "Archived"
	StatusCanceled   = "Canceled"
)

const (
	FieldStatusNotReviewed = "NotReviewed"
	FieldStatusApproved    = "Approved"
	FieldStatusReviewed    = "Reviewed"
	FieldStatusRejected    = "Rejected"
)

var StatusTexts map[string]string = map[string]string{
	StatusIncomplete: "Incomplete",
	StatusReady:      "Ready",
	StatusSubmitted:  "Submitted",
	StatusInReview:   "In Review",
	StatusApproved:   "Approved",
	StatusReviewed:   "Reviewed",
	StatusRejected:   "Rejected",
	StatusArchived:   "Archived",
	StatusCanceled:   "Canceled",
}

var StatusIcons map[string]string = map[string]string{
	StatusIncomplete: "ph:dots-three-outline-fill",
	StatusReady:      "fe:check",
	StatusSubmitted:  "gg:check-o",
	StatusInReview:   "fe:question",
	StatusApproved:   "gg:check-o",
	StatusReviewed:   "fe:warning",
	StatusRejected:   "fe:warning",
	StatusArchived:   "ic:round-lock",
	StatusCanceled:   "fe:outline-close",
}

var StatusColors map[string]string = map[string]string{
	StatusIncomplete: "text-incomplete",
	StatusReady:      "text-ready",
	StatusSubmitted:  "text-submitted",
	StatusInReview:   "text-review",
	StatusApproved:   "text-approved",
	StatusReviewed:   "text-reviewed",
	StatusRejected:   "text-rejected",
	StatusArchived:   "text-archived",
	StatusCanceled:   "text-canceled",
}

type StatusIcon struct {
	Icon     template.URL
	Color    string
	Text     string
	TextSize string
	IconSize int
}

func IsStatusValid(status string) bool {
	_, ok := StatusTexts[status]
	return ok
}

func IsFieldStatusValid(status string) bool {
	switch status {
	case FieldStatusNotReviewed:
		return true
	case FieldStatusApproved:
		return true
	case FieldStatusReviewed:
		return true
	case FieldStatusRejected:
		return true
	default:
		return false
	}
}

type StatusIconParams struct {
	Status      string
	TextSize    string
	IconSize    int
	IncludeText bool
}

func NewStatusIcon(p StatusIconParams) StatusIcon {
	icon, ok := StatusIcons[p.Status]
	if !ok {
		return MakeDefaultStatusIcon(p.IconSize, "text-lg", p.IncludeText)
	}

	color, ok := StatusColors[p.Status]
	if !ok {
		return MakeDefaultStatusIcon(p.IconSize, "text-lg", p.IncludeText)
	}

	result := StatusIcon{
		Icon:     template.URL(icon),
		Color:    color,
		IconSize: p.IconSize,
	}

	if p.IncludeText {
		text, ok := StatusTexts[p.Status]
		if !ok {
			return MakeDefaultStatusIcon(p.IconSize, "text-lg", p.IncludeText)
		}

		result.Text = text
		// TODO: Add some validation to possible size classes
		if len(p.TextSize) == 0 {
			return MakeDefaultStatusIcon(p.IconSize, "text-lg", p.IncludeText)
		}
		result.TextSize = p.TextSize
	}

	return result
}

func MakeDefaultStatusIcon(iconsize int, textsize string, includeText bool) StatusIcon {
	result := StatusIcon{
		Icon:     template.URL(StatusIcons[StatusIncomplete]),
		Color:    StatusColors[StatusIncomplete],
		IconSize: iconsize,
		TextSize: textsize,
	}

	if includeText {
		result.Text = StatusTexts[StatusIncomplete]
	}

	return result
}

func IsEditable(req *query.Request) bool {
	if req.Status == StatusSubmitted {
		return false
	}

	if req.Status == StatusInReview {
		return false
	}

	if req.Status == StatusReviewed {
		return false
	}

	if req.Status == StatusApproved {
		return false
	}

	if req.Status == StatusRejected {
		return false
	}

	if req.Status == StatusCanceled {
		return false
	}

	if req.Status == StatusArchived {
		return false
	}

	return true
}

var ErrNextStatusForbidden error = errors.New("that status update is forbidden")

type NextStatusParams struct {
	Query   *query.Queries
	Request *query.Request
	PID     int64
}

func NextStatus(p NextStatusParams) (string, error) {
	switch p.Request.Status {
	case StatusReady:
		if p.Request.PID != p.PID {
			return "", ErrNextStatusForbidden
		}

		return StatusSubmitted, nil
	case StatusSubmitted:
		if p.Request.PID == p.PID {
			return "", ErrNextStatusForbidden
		}

		return StatusInReview, nil
	case StatusInReview:
		if p.Request.PID == p.PID {
			return "", ErrNextStatusForbidden
		}

		count, err := p.Query.CountCurrentRequestChangeRequestForRequest(context.Background(), p.Request.ID)
		if err != nil {
			return "", err
		}

		if count > 0 {
			return StatusReviewed, nil
		} else {
			return StatusApproved, nil
		}
	case StatusReviewed:
		if p.Request.PID != p.PID {
			return "", ErrNextStatusForbidden
		}

		return StatusReady, nil
	case StatusApproved:
		if p.Request.PID != p.PID {
			return "", ErrNextStatusForbidden
		}

		// TODO: Figure out resolving an approved request
		return "", ErrNextStatusForbidden
	case StatusRejected:
		if p.Request.PID != p.PID {
			return "", ErrNextStatusForbidden
		}

		return StatusArchived, nil
	default:
		return "", ErrNextStatusForbidden
	}
}

type UpdateStatusParams struct {
	Status string
	PID    int64
	RID    int64
}

var (
	ErrInvalidStatus     error = errors.New("invalid status")
	ErrInvalidReviewerID error = errors.New("invalid reviewer ID")
)

func UpdateStatus(q *query.Queries, p UpdateStatusParams) error {
	if !IsStatusValid(p.Status) {
		return ErrInvalidStatus
	}

	if err := q.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     p.RID,
		Status: p.Status,
	}); err != nil {
		return err
	}

	if err := q.CreateHistoryForRequestStatusChange(context.Background(), query.CreateHistoryForRequestStatusChangeParams{
		RID: p.RID,
		PID: p.PID,
	}); err != nil {
		return err
	}

	if p.Status == StatusInReview {
		if p.PID == 0 {
			return ErrInvalidReviewerID
		}

		if err := q.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
			ID:   p.RID,
			RPID: p.PID,
		}); err != nil {
			return err
		}
	}

	if p.Status == StatusReviewed {
		if err := q.LockRequestChangeRequestsForRequest(context.Background(), p.RID); err != nil {
			return err
		}
	}

	return nil
}

type UpdateReadyStatusParams struct {
	Status string
	PID    int64
	RID    int64
	Ready  bool
}

func UpdateReadyStatus(q *query.Queries, p UpdateReadyStatusParams) error {
	if p.Ready && p.Status == StatusIncomplete {
		if err := UpdateStatus(q, UpdateStatusParams{
			PID:    p.PID,
			RID:    p.RID,
			Status: StatusReady,
		}); err != nil {
			return err
		}
	} else if !p.Ready && p.Status == StatusReady {
		if err := UpdateStatus(q, UpdateStatusParams{
			PID:    p.PID,
			RID:    p.RID,
			Status: StatusIncomplete,
		}); err != nil {
			return err
		}
	}
	return nil
}

type CanBePutInReviewParams struct {
	Request             *query.Request
	ReviewerPermissions *player.Permissions
	PID                 int64
}

func CanBePutInReview(p CanBePutInReviewParams) bool {
	if p.PID == p.Request.PID {
		return false
	}

	if p.Request.Status != StatusSubmitted {
		return false
	}

	if !p.ReviewerPermissions.HasPermission(player.PermissionReviewCharacterApplications.Name) {
		return false
	}

	return true
}

type ReviewerTextParams struct {
	Request          *query.Request
	ReviewerUsername string
}

func ReviewerText(p ReviewerTextParams) template.HTML {
	if p.Request.RPID != 0 {
		switch p.Request.Status {
		case StatusInReview:
			var sb strings.Builder
			fmt.Fprintf(&sb, "<span class=\"font-semibold\">Being reviewed by: %s</span>", p.ReviewerUsername)
			return template.HTML(sb.String())
		case StatusApproved:
			var sb strings.Builder
			fmt.Fprintf(&sb, "<span class=\"font-semibold\">Reviewed by: %s</span>", p.ReviewerUsername)
			return template.HTML(sb.String())
		case StatusRejected:
			var sb strings.Builder
			fmt.Fprintf(&sb, "<span class=\"font-semibold\">Reviewed by: %s</span>", p.ReviewerUsername)
			return template.HTML(sb.String())
		case StatusReviewed:
			var sb strings.Builder
			fmt.Fprintf(&sb, "<span class=\"font-semibold\">Reviewed by: %s</span>", p.ReviewerUsername)
			return template.HTML(sb.String())
		default:
			var sb strings.Builder
			fmt.Fprintf(&sb, "<span class=\"font-semibold\">Last reviewed by: %s</span>", p.ReviewerUsername)
			return template.HTML(sb.String())
		}
	} else {
		return template.HTML("<span class=\"font-semibold\">Never reviewed</span>")
	}
}
