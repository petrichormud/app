package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func CharactersPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
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

		b := c.Locals("bind").(fiber.Map)
		b["NewCharacterApplicationPath"] = routes.NewCharacterApplicationPath()
		b["CharacterApplicationSummaries"] = summaries
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render("web/views/characters", b)
	}
}

func CharacterApplicationNamePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		statuses := character.MakeApplicationPartStatuses("name", &app)

		b := c.Locals("bind").(fiber.Map)
		b["Name"] = app.Name
		b["CharacterApplicationNamePath"] = routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		return c.Render("web/views/character/application/flow/name", b, "web/views/layouts/standalone")
	}
}

func CharacterApplicationGenderPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		statuses := character.MakeApplicationPartStatuses("gender", &app)

		gender := character.SanitizeGender(app.Gender)
		b := c.Locals("bind").(fiber.Map)
		b["Name"] = app.Name
		b["GenderNonBinary"] = character.GenderNonBinary
		b["GenderFemale"] = character.GenderFemale
		b["GenderMale"] = character.GenderMale
		b["Gender"] = gender
		b["GenderIsNonBinary"] = gender == character.GenderNonBinary
		b["GenderIsFemale"] = gender == character.GenderFemale
		b["GenderIsMale"] = gender == character.GenderMale
		b["CharacterApplicationGenderPath"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		return c.Render("web/views/character/application/flow/gender", b, "web/views/layouts/standalone")
	}
}

func CharacterApplicationShortDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		statuses := character.MakeApplicationPartStatuses("sdesc", &app)

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["ShortDescription"] = app.ShortDescription
		b["CharacterApplicationShortDescriptionPath"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		return c.Render("web/views/character/application/flow/sdesc", b, "web/views/layouts/standalone")
	}
}

func CharacterApplicationDescriptionPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals(shared.Bind), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		statuses := character.MakeApplicationPartStatuses("description", &app)

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["Description"] = app.Description
		b["CharacterApplicationDescriptionPath"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10))
		b["Statuses"] = statuses
		return c.Render("web/views/character/application/flow/description", b, "web/views/layouts/standalone")
	}
}

func CharacterApplicationBackstoryPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		statuses := character.MakeApplicationPartStatuses("backstory", &app)

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["Backstory"] = app.Backstory
		b["Statuses"] = statuses
		return c.Render("web/views/character/application/flow/backstory", b, "web/views/layouts/standalone")
	}
}

func CharacterApplicationReviewPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		statuses := character.MakeApplicationPartStatuses("review", &app)

		b := c.Locals(shared.Bind).(fiber.Map)
		b["Name"] = app.Name
		b["Statuses"] = statuses
		b["Ready"] = character.IsApplicationReady(&app)
		return c.Render("web/views/character/application/flow/review", b, "web/views/layouts/standalone")
	}
}

func NewCharacterApplication(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		count, err := i.Queries.CountOpenCharacterApplicationsForPlayer(context.Background(), pid.(int64))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if count >= shared.MaxOpenCharacterApplications {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		result, err := i.Queries.CreateRequest(context.Background(), queries.CreateRequestParams{
			Pid:  pid.(int64),
			Type: request.TypeCharacterApplication,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContent(context.Background(), queries.CreateCharacterApplicationContentParams{
			Gender:           character.GenderNonBinary,
			Name:             "",
			ShortDescription: "",
			Description:      "",
			Backstory:        "",
			Rid:              rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		path := routes.CharacterApplicationNamePath(strconv.FormatInt(rid, 10))
		c.Append("HX-Redirect", path)
		return nil
	}
}

func UpdateCharacterApplication(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Name             string `form:"name"`
		Gender           string `form:"gender"`
		ShortDescription string `form:"sdesc"`
		Description      string `form:"description"`
		Backstory        string `form:"backstory"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				req, err := i.Queries.GetRequest(context.Background(), rid)
				if err != nil {
					if err == sql.ErrNoRows {
						c.Status(fiber.StatusNotFound)
						return nil
					}
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				if req.Type != request.TypeCharacterApplication {
					c.Status(fiber.StatusBadRequest)
					return nil
				}
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContent(context.Background(), queries.UpdateCharacterApplicationContentParams{
			Gender:           character.SanitizeGender(r.Gender),
			Name:             r.Name,
			ShortDescription: r.ShortDescription,
			Description:      r.Description,
			Backstory:        r.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationName(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type input struct {
		Name string `form:"name"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		name := character.SanitizeName(r.Name)

		if !character.IsNameValid(name) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		req, err := i.Queries.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Type != request.TypeCharacterApplication {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			Name: name,
			Rid:  rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.CharacterApplicationGenderPath(strconv.FormatInt(rid, 10)))
		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationGender(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type input struct {
		Gender string `form:"gender"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				req, err := i.Queries.GetRequest(context.Background(), rid)
				if err != nil {
					if err == sql.ErrNoRows {
						c.Status(fiber.StatusNotFound)
						return nil
					}
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				if req.Type != request.TypeCharacterApplication {
					c.Status(fiber.StatusBadRequest)
					return nil
				}
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			Gender: r.Gender,
			Rid:    rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		c.Append("HX-Redirect", routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(rid, 10)))
		return nil
	}
}

func UpdateCharacterApplicationShortDescription(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type input struct {
		Sdesc string `form:"sdesc"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				req, err := i.Queries.GetRequest(context.Background(), rid)
				if err != nil {
					if err == sql.ErrNoRows {
						c.Status(fiber.StatusNotFound)
						return nil
					}
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				if req.Type != request.TypeCharacterApplication {
					c.Status(fiber.StatusBadRequest)
					return nil
				}
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentShortDescription(context.Background(), queries.UpdateCharacterApplicationContentShortDescriptionParams{
			ShortDescription: r.Sdesc,
			Rid:              rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		c.Append("HX-Redirect", routes.CharacterApplicationDescriptionPath(strconv.FormatInt(rid, 10)))
		return nil
	}
}

func UpdateCharacterApplicationDescription(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type input struct {
		Description string `form:"description"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				req, err := i.Queries.GetRequest(context.Background(), rid)
				if err != nil {
					if err == sql.ErrNoRows {
						c.Status(fiber.StatusNotFound)
						return nil
					}
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				if req.Type != request.TypeCharacterApplication {
					c.Status(fiber.StatusBadRequest)
					return nil
				}
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			Description: r.Description,
			Rid:         rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.CharacterApplicationBackstoryPath(strconv.FormatInt(rid, 10)))
		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationBackstory(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input for length on the way in
	type input struct {
		Backstory string `form:"backstory"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				req, err := i.Queries.GetRequest(context.Background(), rid)
				if err != nil {
					if err == sql.ErrNoRows {
						c.Status(fiber.StatusNotFound)
						return nil
					}
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				if req.Type != request.TypeCharacterApplication {
					c.Status(fiber.StatusBadRequest)
					return nil
				}
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentBackstory(context.Background(), queries.UpdateCharacterApplicationContentBackstoryParams{
			Backstory: r.Backstory,
			Rid:       rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append("HX-Redirect", routes.CharacterApplicationReviewPath(strconv.FormatInt(rid, 10)))
		c.Status(fiber.StatusOK)
		return nil
	}
}

func UpdateCharacterApplicationVersion(i *shared.Interfaces) fiber.Handler {
	// TODO: Validate this input on the way in
	type input struct {
		Vid int64 `form:"vid"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("web/views/login", c.Locals("bind"), "web/views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		app, err := i.Queries.GetCharacterApplicationContentForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				req, err := i.Queries.GetRequest(context.Background(), rid)
				if err != nil {
					if err == sql.ErrNoRows {
						c.Status(fiber.StatusNotFound)
						return nil
					}
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				if req.Type != request.TypeCharacterApplication {
					c.Status(fiber.StatusBadRequest)
					return nil
				}
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = i.Queries.CreateCharacterApplicationContentHistory(context.Background(), queries.CreateCharacterApplicationContentHistoryParams{
			Gender:           app.Gender,
			Name:             app.Name,
			ShortDescription: app.ShortDescription,
			Description:      app.Description,
			Backstory:        app.Backstory,
			Vid:              app.Vid,
			Rid:              app.Rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		err = i.Queries.UpdateCharacterApplicationContentVersion(context.Background(), queries.UpdateCharacterApplicationContentVersionParams{
			Vid: r.Vid,
			Rid: rid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}
