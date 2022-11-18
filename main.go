package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nitrictech/entgo-planetscale-example/ent"
	"github.com/nitrictech/entgo-planetscale-example/ent/user"
)

var (
	db *ent.Client

	userName      string
	email         string
	userID        int
	userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create a user in the DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			return db.User.Create().
				SetName(userName).
				SetEmail(email).
				Exec(context.TODO())
		},
	}
	userListCmd = &cobra.Command{
		Use:   "list",
		Short: "list the users in the DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			users, err := db.User.Query().All(context.TODO())
			if err != nil {
				return err
			}
			for _, u := range users {
				fmt.Println(u.String())
			}

			return nil
		},
	}
	userDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete the user from the DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := db.User.Delete().Where(user.IDEQ(userID)).Exec(context.TODO())
			return err
		},
	}
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "user DB CRUD commands",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			db, err = mysqlConnectAndMigrate(os.Getenv("DSN"), false)
			return err
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return db.Close()
		},
	}

	migrationExecuteCmd = &cobra.Command{
		Use:   "execute",
		Short: "Execute the migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			db, err = mysqlConnectAndMigrate(os.Getenv("DSN"), true)
			return err
		},
	}
	migrationCreateCmd = &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new migration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createMigration(args[0])
		},
	}
	migrationCmd = &cobra.Command{
		Use: "migration",
	}

	rootCmd = &cobra.Command{
		Use:   "cmd",
		Short: "entgo + planetscale example",
	}
)

func init() {
	// migration commands
	migrationCmd.AddCommand(migrationCreateCmd)
	migrationCmd.AddCommand(migrationExecuteCmd)
	rootCmd.AddCommand(migrationCmd)

	// user CRUD commands
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)
	userCreateCmd.Flags().StringVarP(&userName, "name", "n", "", "-n John Deer")
	userCreateCmd.Flags().StringVarP(&email, "email", "e", "", "-e dearjohn@gmail.com")
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.Flags().IntVarP(&userID, "id", "i", 0, "-i 4")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
