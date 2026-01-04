package bookmarks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a bookmark",
	Long:  "Mark a bookmark as archived.",
	Example: `  clinkding bookmarks archive 42`,
	Args: cobra.ExactArgs(1),
	RunE: runArchive,
}

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive <id>",
	Short: "Unarchive a bookmark",
	Long:  "Mark an archived bookmark as active.",
	Example: `  clinkding bookmarks unarchive 42`,
	Args: cobra.ExactArgs(1),
	RunE: runUnarchive,
}

func runArchive(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	if err := bookmarksAPI.Archive(ctx, id); err != nil {
		return err
	}

	if !cfg.Quiet && !cfg.OutputJSON && !cfg.OutputPlain {
		formatter.Success("Bookmark #%d archived", id)
	}

	return nil
}

func runUnarchive(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	if err := bookmarksAPI.Unarchive(ctx, id); err != nil {
		return err
	}

	if !cfg.Quiet && !cfg.OutputJSON && !cfg.OutputPlain {
		formatter.Success("Bookmark #%d unarchived", id)
	}

	return nil
}
