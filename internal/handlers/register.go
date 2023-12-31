package handlers

import (
	"context"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

func Register(i *shared.Interfaces) fiber.Handler {
	type Player struct {
		Username        string `form:"username"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirmPassword"`
	}

	return func(c *fiber.Ctx) error {
		p := new(Player)

		if err := c.BodyParser(p); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		u := username.Sanitize(p.Username)

		if !username.IsValid(u) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInvalidUsername, layouts.None)
		}

		if !password.IsValid(p.Password) {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInvalidPassword, layouts.None)
		}

		if p.Password != p.ConfirmPassword {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusBadRequest)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInvalidConfirmPassword, layouts.None)
		}

		pwHash, err := password.Hash(p.Password)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)

		}
		defer tx.Rollback()

		qtx := i.Queries.WithTx(tx)
		result, err := qtx.CreatePlayer(
			context.Background(),
			queries.CreatePlayerParams{
				Username: u,
				PwHash:   pwHash,
			},
		)
		if err != nil {
			me, ok := err.(*mysql.MySQLError)
			if !ok {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusInternalServerError)
				return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
			}
			if me.Number == mysqlerr.ER_DUP_ENTRY {
				c.Append("HX-Retarget", "#register-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(shared.HeaderHXAcceptable, "true")
				c.Status(fiber.StatusConflict)
				return c.Render(partials.NoticeSectionError, partials.BindRegisterErrConflict, layouts.None)
			}
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		pid, err := result.LastInsertId()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		err = tx.Commit()
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		username.Cache(i.Redis, pid, p.Username)

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		sess.Set("pid", pid)
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#register-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(shared.HeaderHXAcceptable, "true")
			c.Status(fiber.StatusInternalServerError)
			return c.Render(partials.NoticeSectionError, partials.BindRegisterErrInternal, layouts.None)
		}

		c.Append("HX-Refresh", "true")
		c.Status(fiber.StatusCreated)
		return nil
	}
}
