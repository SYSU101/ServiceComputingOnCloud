package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:     "agenda",
	Version: "0.0.1",
}

func init() {
	rootCmd.AddCommand(RegisterCmd(), LoginCmd(), LogoutCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
