package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func RecoverPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(view.Recover, view.Bind(c), layout.Standalone)
	}
}

func SearchPlayer(i *service.Interfaces) fiber.Handler {
	type input struct {
		Search string `form:"search"`
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		searchStr := fmt.Sprintf("%%%s%%", r.Search)
		players, err := i.Queries.SearchPlayersByUsername(context.Background(), searchStr)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		dest := c.Params("dest")
		if len(dest) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if dest == "player-permissions" {
			// TODO: Move this to a constant and inject it
			b := view.Bind(c)
			b["Players"] = players
			c.Status(fiber.StatusOK)
			return c.Render(partial.PlayerPermissionsSearchResults, b, layout.None)
		}

		c.Status(fiber.StatusBadRequest)
		return nil
	}
}

func PlayerPermissionsPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		perms, ok := lperms.(player.Permissions)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermissionInSet(player.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		return c.Render(view.PlayerPermissions, view.Bind(c), layout.Main)
	}
}

func PlayerPermissionsDetailPage(i *service.Interfaces) fiber.Handler {
	type playerPermissionDetail struct {
		Permission string
		About      string
		Link       string
		Granted    bool
		Disabled   bool
	}
	return func(c *fiber.Ctx) error {
		pid := c.Locals("pid")

		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		lperms := c.Locals("perms")
		if lperms == nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		iperms, ok := lperms.(player.Permissions)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if !iperms.HasPermissionInSet(player.ShowPermissionViewPermissions) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		u := c.Params("username")
		if len(u) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), u)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		pperms, err := qtx.ListPlayerPermissions(context.Background(), p.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		perms := player.NewPermissions(p.ID, pperms)
		allPerms := []fiber.Map{}
		for _, perm := range player.AllPermissions {
			granted := perms.Permissions[perm.Name]
			disabled := true
			if granted {
				disabled = !iperms.CanRevokePermission(perm.Name)
			} else {
				disabled = !iperms.CanGrantPermission(perm.Name)
			}
			// TODO: Fix this to just use the "Name" for the permission
			pm := fiber.Map{
				"Name":     perm.Name,
				"Tag":      perm.Name,
				"Title":    perm.Title,
				"About":    perm.About,
				"Link":     route.PlayerPermissionsTogglePath(strconv.FormatInt(p.ID, 10), perm.Name),
				"Granted":  granted,
				"Disabled": disabled,
			}
			allPerms = append(allPerms, pm)
		}

		b := view.Bind(c)
		b["Username"] = u
		b["Permissions"] = allPerms
		return c.Render(view.PlayerPermissionsDetail, b)
	}
}

func TogglePlayerPermission(i *service.Interfaces) fiber.Handler {
	type input struct {
		Grant bool `form:"issued"`
	}
	return func(c *fiber.Ctx) error {
		ipid := c.Locals("pid")
		if ipid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !perms.HasPermission(player.PermissionGrantAll.Name) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		ppid := c.Params("id")
		if len(ppid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		pid, err := strconv.ParseInt(ppid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		ptag := c.Params("tag")
		if len(ptag) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !player.IsValidPermissionName(ptag) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		pperms, err := qtx.ListPlayerPermissions(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		r := new(input)
		if err = c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		perms = player.NewPermissions(pid, pperms)
		perm := player.AllPermissionsByName[ptag]
		_, granted := perms.Permissions[perm.Name]

		if r.Grant && granted {
			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			c.Status(fiber.StatusConflict)
			return nil
		}

		if r.Grant && !granted {
			if !perms.CanGrantPermission(perm.Name) {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			params := query.CreatePlayerPermissionIssuedChangeHistoryParams{
				IPID: ipid.(int64),
				PID:  pid,
				Name: perm.Name,
			}
			if err := qtx.CreatePlayerPermissionIssuedChangeHistory(context.Background(), params); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			_, err := qtx.CreatePlayerPermission(context.Background(), query.CreatePlayerPermissionParams{
				IPID: ipid.(int64),
				PID:  pid,
				Name: perm.Name,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			return nil
		}

		if !r.Grant && granted {
			if !perms.CanRevokePermission(perm.Name) {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			params := query.CreatePlayerPermissionRevokedChangeHistoryParams{
				IPID: ipid.(int64),
				PID:  pid,
				Name: perm.Name,
			}
			if err := qtx.CreatePlayerPermissionRevokedChangeHistory(context.Background(), params); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err := qtx.DeletePlayerPermission(context.Background(), query.DeletePlayerPermissionParams{
				PID:  pid,
				Name: perm.Name,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			return nil
		}

		if !r.Grant && !granted {
			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}
}
