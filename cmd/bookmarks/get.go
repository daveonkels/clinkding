package bookmarks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a bookmark by ID",
	Long:  "Retrieve detailed information about a specific bookmark.",
	Example: `  clinkding bookmarks get 42
  clinkding bookmarks get 42 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func runGet(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	bookmark, err := bookmarksAPI.Get(ctx, id)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(bookmark)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(
			strconv.Itoa(bookmark.ID),
			bookmark.URL,
			bookmark.Title,
			bookmark.Description,
			bookmark.Notes,
			strings.Join(bookmark.TagNames, ","),
		)
		return nil
	}

	// Human-friendly output
	formatter.Println(formatter.Bold("Bookmark #%d"), bookmark.ID)
	formatter.Println("")
	formatter.Println("URL:         %s", bookmark.URL)
	formatter.Println("Title:       %s", bookmark.Title)
	if bookmark.Description != "" {
		formatter.Println("Description: %s", bookmark.Description)
	}
	if bookmark.Notes != "" {
		formatter.Println("Notes:       %s", bookmark.Notes)
	}
	if len(bookmark.TagNames) > 0 {
		formatter.Println("Tags:        %s", strings.Join(bookmark.TagNames, ", "))
	}
	formatter.Println("Archived:    %t", bookmark.IsArchived)
	formatter.Println("Unread:      %t", bookmark.Unread)
	formatter.Println("Shared:      %t", bookmark.Shared)
	formatter.Println("Added:       %s", bookmark.DateAdded.Format(time.RFC3339))
	formatter.Println("Modified:    %s", bookmark.DateModified.Format(time.RFC3339))

	return nil
}
