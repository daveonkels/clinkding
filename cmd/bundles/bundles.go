package bundles

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "bundles",
	Short: "Manage bundles",
	Long:  "Commands for listing, creating, updating, and deleting bundles.",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)
}
