package assets

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "assets",
	Short: "Manage bookmark assets",
	Long:  "Commands for listing, uploading, downloading, and deleting bookmark assets (file attachments).",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(uploadCmd)
	Cmd.AddCommand(downloadCmd)
	Cmd.AddCommand(deleteCmd)
}
