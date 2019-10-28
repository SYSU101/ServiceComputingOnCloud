package cmd

import (
	"github.com/SYSU101/agenda/entity"
	"github.com/SYSU101/agenda/logger"
	"github.com/spf13/cobra"
)

func LogoutCmd() *cobra.Command {
	logoutCmd := &cobra.Command{
		Use:   "logout <options>",
		Short: "Log out",
		RunE: func(_ *cobra.Command, _ []string) error {
			logger.Printf("[info]  logging out\n")
			return entity.Logout()
		},
		PostRunE: func(_ *cobra.Command, _ []string) error {
			return entity.SaveUser()
		},
	}
	return logoutCmd
}
