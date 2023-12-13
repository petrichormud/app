package handlers

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func CharactersPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals(shared.Bind), "web/views/layouts/standalone")
		}

		apps, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("web/views/500", c.Locals(shared.Bind))
		}

		summaries := []character.ApplicationSummary{}
		for _, app := range apps {
			summaries = append(summaries, character.NewSummaryFromApplication(&app.Request, &app.CharacterApplicationContent))
		}

		b := c.Locals(shared.Bind).(fiber.Map)
		b["NewCharacterApplicationPath"] = routes.NewCharacterApplicationPath()
		b["CharacterApplicationSummaries"] = summaries
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render("web/views/characters", b)
	}
}
