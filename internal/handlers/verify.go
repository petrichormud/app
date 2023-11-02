package handlers

import (
	"context"
	"slices"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/username"
)

func Verify(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("t")
		exists, err := i.Redis.Exists(context.Background(), token).Result()
		if err != nil {
			return c.Redirect("/")
		}
		if exists != 1 {
			return c.Redirect("/")
		}

		pid := c.Locals("pid")
		if pid == nil {
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		eid, err := i.Redis.Get(context.Background(), token).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			return c.Redirect("/")
		}
		e, err := i.Queries.GetEmail(context.Background(), id)
		if err != nil {
			// TODO: Distinguish between "not found" and a connection error
			return c.Redirect("/")
		}

		perms, err := permissions.List(i, pid.(int64))
		if err != nil {
			return c.Redirect("/")
		}
		if !slices.Contains(perms, permissions.AddEmail) {
			return c.Redirect("/")
		}

		un, err := username.Get(i.Redis, pid.(int64))
		if err != nil {
			return c.Redirect("/")
		}
		b := c.Locals("bind").(fiber.Map)
		b["VerifyToken"] = c.Query("t")
		b["Address"] = e.Address
		b["Username"] = un

		return c.Render("web/views/verify", b, "web/views/layouts/standalone")
	}
}

func VerifyEmail(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			// TODO: This should redirect them back to the login page for this token
			return nil
		}

		perms, err := permissions.List(i, pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if !slices.Contains(perms, permissions.AddEmail) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		key := c.Query("t")
		if len(key) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		eid, err := i.Redis.Get(context.Background(), key).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		id, err := strconv.ParseInt(eid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.MarkEmailVerified(context.Background(), id)
		if err != nil {
			// TODO: Distinguish between a "not found" error and a connection error
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Redis.Del(context.Background(), key).Err()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return c.Render("web/views/partials/verify/success", &fiber.Map{}, "")
	}
}
