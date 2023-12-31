package request

import (
	"fmt"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

// TODO: This is a pure dump from the character module; this badly needs cleaning up

type ReviewDialogData struct {
	Path     string
	Variable string
}

type ApplicationSummary struct {
	StatusIcon       StatusIcon
	ReviewDialog     ReviewDialogData
	Status           string
	StatusColor      string
	StatusText       string
	Link             string
	Name             string
	Author           string
	Reviewer         string
	ID               int64
	RPID             int64
	StatusIncomplete bool
	StatusReady      bool
	StatusSubmitted  bool
	StatusInReview   bool
	StatusApproved   bool
	StatusReviewed   bool
	StatusRejected   bool
	StatusArchived   bool
	StatusCanceled   bool
	Reviewed         bool
}

func NewSummaryFromApplication(p *queries.Player, reviewer string, req *queries.Request, app *queries.CharacterApplicationContent) ApplicationSummary {
	name := app.Name
	if len(app.Name) == 0 {
		name = constants.DefaultName
	}

	reviewed := len(reviewer) > 0

	return ApplicationSummary{
		Reviewed: reviewed,
		ReviewDialog: ReviewDialogData{
			Path:     routes.RequestStatusPath(req.ID),
			Variable: fmt.Sprintf("showReviewDialogFor%s%s", app.Name, p.Username),
		},
		Status:           req.Status,
		StatusColor:      StatusColors[req.Status],
		StatusText:       StatusTexts[req.Status],
		StatusIncomplete: req.Status == StatusIncomplete,
		StatusReady:      req.Status == StatusReady,
		StatusSubmitted:  req.Status == StatusSubmitted,
		StatusInReview:   req.Status == StatusInReview,
		StatusApproved:   req.Status == StatusApproved,
		StatusReviewed:   req.Status == StatusReviewed,
		StatusRejected:   req.Status == StatusRejected,
		StatusArchived:   req.Status == StatusArchived,
		StatusCanceled:   req.Status == StatusCanceled,
		Link:             routes.RequestPath(req.ID),
		ID:               req.ID,
		Name:             name,
		Author:           p.Username,
		Reviewer:         reviewer,
		StatusIcon:       MakeStatusIcon(MakeStatusIconParams{Status: req.Status, Size: "48", IncludeText: false}),
	}
}
