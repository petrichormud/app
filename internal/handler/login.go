package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player/password"
	"petrichormud.com/app/internal/player/username"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/view"
)

func Login(i *service.Interfaces) fiber.Handler {
	type input struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), in.Username)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		v, err := password.Verify(in.Password, p.PwHash)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}
		if !v {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		pid := p.ID
		err = username.Cache(i.Redis, pid, p.Username)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		settings, err := qtx.GetPlayerSettings(context.Background(), pid)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This means a player got created without settings
				c.Status(fiber.StatusInternalServerError)
				c.Append("HX-Retarget", "#login-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusUnauthorized)
				return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		sess, err := i.Sessions.Get(c)
		if err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		sess.Set("pid", pid)
		theme := sess.Get("theme")
		if theme != nil {
			if err := qtx.UpdatePlayerSettingsTheme(context.Background(), query.UpdatePlayerSettingsThemeParams{
				PID:   pid,
				Theme: theme.(string),
			}); err != nil {
				c.Append("HX-Retarget", "#login-error")
				c.Append("HX-Reswap", "outerHTML")
				c.Append(header.HXAcceptable, "true")
				c.Status(fiber.StatusUnauthorized)
				return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
			}
		} else {
			sess.Set("theme", settings.Theme)
		}
		if err = sess.Save(); err != nil {
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append("HX-Retarget", "#login-error")
			c.Append("HX-Reswap", "outerHTML")
			c.Append(header.HXAcceptable, "true")
			c.Status(fiber.StatusUnauthorized)
			return c.Render(partial.NoticeSectionError, partial.BindLoginErr, layout.None)
		}

		c.Append("HX-Refresh", "true")
		return nil
	}
}

func LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Rework this to use the util
		pid := c.Locals("pid")
		if pid != nil {
			return c.Redirect(route.Home)
		}

		return c.Render(view.Login, view.Bind(c), layout.Standalone)
	}
}
