package bookmarks

import (
	"context"
	"fmt"
	"strings"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <url>",
	Short: "Check if a URL exists",
	Long:  "Check if a URL already exists in your bookmarks and get scraped metadata.",
	Example: `  clinkding bookmarks check https://go.dev
  clinkding bookmarks check https://go.dev --json`,
	Args: cobra.ExactArgs(1),
	RunE: runCheck,
}

func runCheck(cobraCmd *cobra.Command, args []string) error {
	urlToCheck := args[0]

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	result, err := bookmarksAPI.Check(ctx, urlToCheck)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(result)
	}

	if cfg.OutputPlain {
		if result.Bookmark != nil {
			output.PrintPlainLine("exists", fmt.Sprintf("%d", result.Bookmark.ID))
		} else {
			output.PrintPlainLine("new", result.Metadata.Title, result.Metadata.Description)
		}
		return nil
	}

	// Human-friendly output
	if result.Bookmark != nil {
		formatter.Warning("Bookmark already exists!")
		formatter.Println("")
		formatter.Println("ID:    %d", result.Bookmark.ID)
		formatter.Println("Title: %s", result.Bookmark.Title)
		formatter.Println("Tags:  %s", strings.Join(result.Bookmark.TagNames, ", "))
	} else {
		formatter.Success("URL not bookmarked yet")
		formatter.Println("")
		formatter.Println("Scraped metadata:")
		formatter.Println("  Title:       %s", result.Metadata.Title)
		formatter.Println("  Description: %s", result.Metadata.Description)
		if len(result.AutoTags) > 0 {
			formatter.Println("  Suggested tags: %s", strings.Join(result.AutoTags, ", "))
		}
		formatter.Println("")
		formatter.Info("Create with: clinkding bookmarks create %s", urlToCheck)
	}

	return nil
}
