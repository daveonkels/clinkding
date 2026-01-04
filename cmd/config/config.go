package config

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage clinkding configuration",
	Long:  "Commands for initializing, viewing, and testing clinkding configuration.",
}

func init() {
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(testCmd)
}
