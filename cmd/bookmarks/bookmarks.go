package bookmarks

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "Manage bookmarks",
	Long:  "Commands for listing, creating, updating, and deleting bookmarks.",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(checkCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(archiveCmd)
	Cmd.AddCommand(unarchiveCmd)
	Cmd.AddCommand(deleteCmd)
}
