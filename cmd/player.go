/*
Copyright © 2023 Alec DuBois <alec@petrichormud.com>
*/
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/player/password"
	"petrichormud.com/app/internal/player/username"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/service"
)

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "Add, edit, and get information about players.",
}

var addPlayerCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new player.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		u, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		pw, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		dbURL, err := cmd.Flags().GetString("db-url")
		if err != nil {
			return err
		}

		if !username.IsValid(u) {
			return errors.New("please enter a valid username")
		}

		if !password.IsValid(pw) {
			return errors.New("please enter a valid password")
		}

		db, err := sql.Open("mysql", dbURL)
		if err != nil {
			return err
		}
		if err = service.SetupDB(db); err != nil {
			return errors.New("error while setting up DB")
		}
		if err = service.PingDB(db); err != nil {
			return errors.New("error while pinging DB")
		}

		q := query.New(db)
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		qtx := q.WithTx(tx)

		pwHash, err := password.Hash(pw)
		if err != nil {
			return errors.New("error while hashing password")
		}

		result, err := qtx.CreatePlayer(context.Background(), query.CreatePlayerParams{
			Username: u,
			PwHash:   pwHash,
		})
		if err != nil {
			return err
		}

		pid, err := result.LastInsertId()
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		msg := fmt.Sprintf("User %s created with PID %d.", u, pid)
		fmt.Println(msg)
		return nil
	},
}

var playerPermissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "Grant, revoke, or get data about player permissions.",
}

var grantPlayerPermissionCmd = &cobra.Command{
	Use:   "grant",
	Short: "Grant a permission to a player.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		dbURL, err := cmd.Flags().GetString("db-url")
		if err != nil {
			return err
		}
		u, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		ptag, err := cmd.Flags().GetString("permission")
		if err != nil {
			return err
		}

		if !username.IsValid(u) {
			return errors.New("please enter a valid username")
		}

		perm, ok := player.AllPermissionsByName[ptag]
		if !ok {
			return errors.New("please enter a valid permission tag")
		}

		db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", dbURL))
		if err != nil {
			return err
		}
		if err = service.SetupDB(db); err != nil {
			return errors.New("error while setting up DB")
		}
		if err = service.PingDB(db); err != nil {
			return errors.New("error while pinging DB")
		}

		q := query.New(db)
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		qtx := q.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), u)
		if err != nil {
			return err
		}

		ps, err := qtx.ListPlayerPermissions(context.Background(), p.ID)
		if err != nil {
			return err
		}

		perms := player.NewPermissions(p.ID, ps)
		_, granted := perms.Permissions[perm.Name]
		if granted {
			msg := fmt.Sprintf("Player %s already has permission %s.", u, perm.Name)
			fmt.Println(msg)
			return nil
		}

		if err := qtx.CreatePlayerPermissionIssuedChangeHistory(context.Background(), query.CreatePlayerPermissionIssuedChangeHistoryParams{
			PID:  p.ID,
			IPID: p.ID,
			Name: perm.Name,
		}); err != nil {
			return err
		}
		_, err = qtx.CreatePlayerPermission(context.Background(), query.CreatePlayerPermissionParams{
			PID:  p.ID,
			IPID: p.ID,
			Name: perm.Name,
		})
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		msg := fmt.Sprintf("User %s granted permission %s.", u, perm.Name)
		fmt.Println(msg)
		return nil
	},
}

// TODO: Split out the username case here into a "granted" command
var listPlayerPermissionCmd = &cobra.Command{
	Use:   "list",
	Short: "List a player's current permissions.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		dbURL, err := cmd.Flags().GetString("db-url")
		if err != nil {
			return err
		}
		u, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}

		if !username.IsValid(u) {
			permissionnames := []string{}
			for _, permission := range player.AllPermissions {
				permissionnames = append(permissionnames, permission.Name)
			}
			permmissionsmsg := strings.Join(permissionnames, "\n")
			msg := fmt.Sprintf("All permissions:\n%s", permmissionsmsg)
			fmt.Println(msg)
			return nil
		}

		db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", dbURL))
		if err != nil {
			return err
		}
		if err = service.SetupDB(db); err != nil {
			return errors.New("error while setting up DB")
		}
		if err = service.PingDB(db); err != nil {
			return errors.New("error while pinging DB")
		}

		q := query.New(db)
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		qtx := q.WithTx(tx)

		p, err := qtx.GetPlayerByUsername(context.Background(), u)
		if err != nil {
			return err
		}

		ps, err := qtx.ListPlayerPermissions(context.Background(), p.ID)
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		perms := player.NewPermissions(p.ID, ps)
		msg := fmt.Sprintf("User %s has permissions %s.", u, strings.Join(perms.PermissionsList, ", "))
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)

	playerCmd.AddCommand(addPlayerCmd)
	addPlayerCmd.Flags().StringP("db-url", "d", "root:pass@/test", "The URL for the database.")
	addPlayerCmd.Flags().StringP("username", "u", "", "The username for the new player.")
	addPlayerCmd.Flags().StringP("password", "p", "", "The password for the player.")
	addPlayerCmd.MarkFlagRequired("username")
	addPlayerCmd.MarkFlagRequired("password")

	playerCmd.AddCommand(playerPermissionCmd)

	playerPermissionCmd.AddCommand(grantPlayerPermissionCmd)
	grantPlayerPermissionCmd.Flags().StringP("username", "u", "", "The username for the player.")
	grantPlayerPermissionCmd.Flags().StringP("permission", "p", "", "The tag for the permission to grant.")
	grantPlayerPermissionCmd.Flags().StringP("db-url", "d", "root:pass@/test", "The URL for the database.")

	playerPermissionCmd.AddCommand(listPlayerPermissionCmd)
	listPlayerPermissionCmd.Flags().StringP("username", "u", "", "The username for the player.")
	listPlayerPermissionCmd.Flags().StringP("db-url", "d", "root:pass@/test", "The URL for the database.")
}
