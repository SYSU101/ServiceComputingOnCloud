package cmd

import (
	"github.com/SYSU101/agenda/entity"
	"github.com/SYSU101/agenda/fmterror"
	"github.com/SYSU101/agenda/logger"
	"github.com/spf13/cobra"
)

func RegisterCmd() *cobra.Command {
	registerUser := &entity.User{
		Username: "",
		Password: "",
		Email:    "",
	}
	registerCmd := &cobra.Command{
		Use:   "register <options>",
		Short: "Register a new user",
		RunE: func(_ *cobra.Command, _ []string) error {
			logger.Printf("[info]  regsitering new user with\n\tusername: %v\n\temail: %v\n", registerUser.Username, registerUser.Email)
			if entity.CurrentUser != nil {
				return fmterror.New("There has been other user logged in, pleas log out first")
			} else {
				return entity.Register(registerUser)
			}
		},
		PostRunE: func(_ *cobra.Command, _ []string) error {
			return entity.SaveUser()
		},
	}
	registerCmd.Flags().StringVarP(&registerUser.Username, "username", "u", "", "username of new user")
	registerCmd.Flags().StringVarP(&registerUser.Password, "password", "p", "", "(optional)password of new user")
	registerCmd.Flags().StringVarP(&registerUser.Email, "email", "e", "", "email of new user")
	return registerCmd
}
