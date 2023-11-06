package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/email"
	"petrichormud.com/app/internal/shared"
)

const ResendRoute = "/player/email/:id/resend"

func Resend(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			// TODO: Check that this behavior is accurate - maybe send back a 403 instead
			// TODO: Along with a minor "log in again" component
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Append("HX-Refresh", "true")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		eid := c.Params("id")
		if len(eid) == 0 {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{}, "")
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{}, "")
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{}, "")
		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)

		e, err := qtx.GetEmail(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusNotFound)
				return c.Render("web/views/partials/profile/email/resend/err-404", &fiber.Map{
					"ID": id,
				}, "")
			}
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{
				"ID": id,
			}, "")
		}

		if e.Verified {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusConflict)
			return c.Render("web/views/partials/profile/email/resend/err-conflict", &fiber.Map{}, "")
		}
		if e.Pid != pid.(int64) {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusForbidden)
			// TODO: Make this a different error - here, the caller doesn't own the email
			return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{}, "")
		}

		_, err = qtx.GetVerifiedEmailByAddress(context.Background(), e.Address)
		if err != nil {
			if err != sql.ErrNoRows {
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
				return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{}, "")
			}
		}
		if err == nil {
			// TODO: This is a new error state - it means another user has claimed and verified the email before you
			c.Status(fiber.StatusForbidden)
			return nil
		}

		err = email.Verify(i.Redis, id, e.Address)
		if err != nil {
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/partials/profile/email/resend/err-internal", &fiber.Map{}, "")
		}

		return c.Render("web/views/partials/profile/email/resend/success", &fiber.Map{}, "")
	}
}