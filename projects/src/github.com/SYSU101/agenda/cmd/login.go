package cmd

import (
	"github.com/SYSU101/agenda/entity"
	"github.com/SYSU101/agenda/logger"
	"github.com/spf13/cobra"
)

func LoginCmd() *cobra.Command {
	username, password := "", ""
	loginCmd := &cobra.Command{
		Use:   "login <options>",
		Short: "Log in",
		RunE: func(_ *cobra.Command, _ []string) error {
			logger.Printf("[info]  logging in with\n\tusername: %v\n", username)
			return entity.Login(username, password)
		},
		PostRunE: func(_ *cobra.Command, _ []string) error {
			return entity.SaveUser()
		},
	}
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "username of user")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "(optional)password of user")
	return loginCmd
}
