package user

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "user",
	Short: "User profile and settings",
	Long:  "Commands for viewing user profile and settings.",
}

func init() {
	Cmd.AddCommand(profileCmd)
}
