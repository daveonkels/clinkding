package main

import (
	"fmt"
	"os"

	"github.com/daveonkels/clinkding/cmd"
	assetsCmd "github.com/daveonkels/clinkding/cmd/assets"
	bookmarksCmd "github.com/daveonkels/clinkding/cmd/bookmarks"
	bundlesCmd "github.com/daveonkels/clinkding/cmd/bundles"
	configCmd "github.com/daveonkels/clinkding/cmd/config"
	tagsCmd "github.com/daveonkels/clinkding/cmd/tags"
	userCmd "github.com/daveonkels/clinkding/cmd/user"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)

	// Register commands
	cmd.AddCommand(assetsCmd.Cmd)
	cmd.AddCommand(bookmarksCmd.Cmd)
	cmd.AddCommand(bundlesCmd.Cmd)
	cmd.AddCommand(configCmd.Cmd)
	cmd.AddCommand(tagsCmd.Cmd)
	cmd.AddCommand(userCmd.Cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
