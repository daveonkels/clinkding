package tags

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "tags",
	Short: "Manage tags",
	Long:  "Commands for listing, getting, and creating tags.",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(createCmd)
}
