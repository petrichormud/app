package request

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/definition"
	"petrichormud.com/app/internal/request/dialog"
)

var DialogsByType map[string]*dialog.DefinitionGroup = map[string]*dialog.DefinitionGroup{
	TypeCharacterApplication: &definition.DialogsCharacterApplication,
}

const (
	BindCancelDialog       = "CancelDialog"
	BindSubmitDialog       = "SubmitDialog"
	BindPutInReviewDialog  = "PutInReviewDialog"
	BindApproveDialog      = "ApproveDialog"
	BindFinishReviewDialog = "FinishReviewDialog"
	BindRejectDialog       = "RejectDialog"
)

// TODO: Move Dialogs to new struct
func BindDialogs(b fiber.Map, req *query.Request) (fiber.Map, error) {
	dialogs, ok := DialogsByType[req.Type]
	if !ok {
		// TODO: Discrete error
		return fiber.Map{}, ErrNoDefinition
	}
	dialogs.SetPath(req.ID)

	b["Dialogs"] = dialogs

	b[BindCancelDialog] = dialogs.Cancel
	b[BindSubmitDialog] = dialogs.Submit
	b[BindPutInReviewDialog] = dialogs.PutInReview
	b[BindApproveDialog] = dialogs.Approve
	b[BindFinishReviewDialog] = dialogs.FinishReview

	return b, nil
}
