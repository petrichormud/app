package handlers

import (
	"context"
	"database/sql"
	"html/template"
	"slices"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/views"
)

func HelpPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		headers, err := qtx.ListHelpTitleAndSub(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		help := []fiber.Map{}
		for _, header := range headers {
			help = append(help, fiber.Map{
				"Title": header.Title,
				"Sub":   header.Sub,
				"Path":  routes.HelpFilePath(header.Slug),
			})
		}

		b := views.Bind(c)
		b["Help"] = help
		return c.Render(views.Help, b, layouts.Main)
	}
}

func HelpFilePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		slugs, err := qtx.ListHelpSlugs(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		slug := c.Params("slug")
		if !slices.Contains(slugs, slug) {
			c.Status(fiber.StatusNotFound)
			return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
		}

		help, err := qtx.GetHelp(context.Background(), slug)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		relatedRecords, err := qtx.GetHelpRelated(context.Background(), slug)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		related := []fiber.Map{}
		for _, record := range relatedRecords {
			related = append(related, fiber.Map{
				"Title": record.RelatedTitle,
				"Sub":   record.RelatedSub,
				"Path":  routes.HelpFilePath(record.Related),
			})
		}

		html := template.HTML(help.HTML)
		b := views.Bind(c)
		b["Content"] = html
		b["Related"] = related
		// TODO: Once the help path can take a query string, save the last state of the session's help path
		b["HelpPath"] = routes.Help
		return c.Render(views.HelpFile, b, layouts.Main)
	}
}