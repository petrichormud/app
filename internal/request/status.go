package request

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/field"
	"petrichormud.com/app/internal/request/status"
)

const (
	StatusIncomplete = status.Incomplete
	StatusReady      = status.Ready
	StatusSubmitted  = status.Submitted
	StatusInReview   = status.InReview
	StatusApproved   = status.Approved
	StatusReviewed   = status.Reviewed
	StatusRejected   = status.Rejected
	StatusFulfilled  = status.Fulfilled
	StatusArchived   = status.Archived
	StatusCanceled   = status.Canceled
)

const (
	FieldStatusNotReviewed = field.StatusNotReviewed
	FieldStatusApproved    = field.StatusApproved
	FieldStatusReviewed    = field.StatusReviewed
	FieldStatusRejected    = field.StatusRejected
)

// TODO: Move the rest of this file to the underlying status package?
var StatusTexts map[string]string = map[string]string{
	StatusIncomplete: "Incomplete",
	StatusReady:      "Ready",
	StatusSubmitted:  "Submitted",
	StatusInReview:   "In Review",
	StatusApproved:   "Approved",
	StatusReviewed:   "Reviewed",
	StatusRejected:   "Rejected",
	StatusFulfilled:  "Fulfilled",
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
	StatusFulfilled:  "gg:check-o",
	StatusArchived:   "ic:round-lock",
	StatusCanceled:   "fe:outline-close",
}

// TODO: Get a color value for text-fulfilled
var StatusColors map[string]string = map[string]string{
	StatusIncomplete: "text-incomplete",
	StatusReady:      "text-ready",
	StatusSubmitted:  "text-submitted",
	StatusInReview:   "text-review",
	StatusApproved:   "text-approved",
	StatusReviewed:   "text-reviewed",
	StatusRejected:   "text-rejected",
	StatusFulfilled:  "text-submitted",
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

// TODO: Split these out into better errors
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

		// TODO: Do this by the fieldmap instead of counting the change requests
		count, err := p.Query.CountOpenRequestChangeRequestsForRequest(context.Background(), p.Request.ID)
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

		return StatusSubmitted, nil
	case StatusApproved:
		fulfilledby, err := FulfilledBy(p.Request.Type)
		if err != nil {
			return "", err
		}

		if fulfilledby == "Player" && p.Request.PID == p.PID {
			return StatusFulfilled, nil
		}

		if fulfilledby == "Reviewer" && p.Request.RPID == p.PID {
			return StatusFulfilled, nil
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
	Request *query.Request
	Status  string
	PID     int64
}

var (
	ErrInvalidStatus     error = errors.New("invalid status")
	ErrInvalidReviewerID error = errors.New("invalid reviewer ID")
)

// TODO: Invalidate an invalid status? i.e., make this a state machine?
func UpdateStatus(q *query.Queries, p UpdateStatusParams) error {
	if !IsStatusValid(p.Status) {
		return ErrInvalidStatus
	}

	if p.Status == StatusInReview {
		if p.PID == 0 {
			return ErrInvalidReviewerID
		}

		if err := q.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
			ID:   p.Request.ID,
			RPID: p.PID,
		}); err != nil {
			return err
		}
	}

	if p.Status == StatusFulfilled {
		// TODO: Validate that the correct person is attempting to fulfill this?
		fulfilledby, err := FulfilledBy(p.Request.Type)
		if err != nil {
			return err
		}
		// TODO: Use a constant here instead
		if fulfilledby == "Player" && p.PID != p.Request.PID {
			return ErrNextStatusForbidden
		}
		// TODO: Use a constant here instead
		if fulfilledby == "Reviewer" && p.PID != p.Request.RPID {
			return ErrNextStatusForbidden
		}
	}

	if err := q.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     p.Request.ID,
		Status: p.Status,
	}); err != nil {
		return err
	}

	return nil
}

type CanBePutInReviewParams struct {
	Request     *query.Request
	Permissions *player.Permissions
	PID         int64
}

func CanBePutInReview(p CanBePutInReviewParams) bool {
	if p.PID == p.Request.PID {
		return false
	}
	if p.Request.Status != StatusSubmitted {
		return false
	}
	if !p.Permissions.HasPermission(player.PermissionReviewCharacterApplications.Name) {
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
