package bookmarks

import (
	"context"
	"strconv"
	"strings"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var (
	createTitle       string
	createDescription string
	createNotes       string
	createTags        string
	createShared      bool
	createUnread      bool
	createNoScraping  bool
)

var createCmd = &cobra.Command{
	Use:   "create <url>",
	Short: "Create a new bookmark",
	Long:  "Add a new bookmark with optional metadata.",
	Example: `  clinkding bookmarks create https://go.dev
  clinkding bookmarks create https://go.dev --title "Go" --tags "golang,programming"
  clinkding bookmarks create https://go.dev --no-scraping --unread`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVar(&createTitle, "title", "", "bookmark title")
	createCmd.Flags().StringVar(&createDescription, "description", "", "description")
	createCmd.Flags().StringVar(&createNotes, "notes", "", "notes")
	createCmd.Flags().StringVar(&createTags, "tags", "", "comma-separated tags")
	createCmd.Flags().BoolVar(&createShared, "shared", false, "make bookmark shared")
	createCmd.Flags().BoolVar(&createUnread, "unread", false, "mark as unread")
	createCmd.Flags().BoolVar(&createNoScraping, "no-scraping", false, "skip automatic metadata scraping")
}

func runCreate(cobraCmd *cobra.Command, args []string) error {
	urlToCreate := args[0]

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	// Parse tags
	var tags []string
	if createTags != "" {
		tags = strings.Split(createTags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	// Create bookmark request
	bookmarkCreate := &models.BookmarkCreate{
		URL:         urlToCreate,
		Title:       createTitle,
		Description: createDescription,
		Notes:       createNotes,
		TagNames:    tags,
		Shared:      createShared,
		Unread:      createUnread,
	}

	ctx := context.Background()
	bookmark, err := bookmarksAPI.Create(ctx, bookmarkCreate)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(bookmark)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(strconv.Itoa(bookmark.ID), bookmark.URL, bookmark.Title)
		return nil
	}

	// Human-friendly output
	formatter.Success("Bookmark created!")
	formatter.Println("")
	formatter.Println("ID:    %d", bookmark.ID)
	formatter.Println("Title: %s", bookmark.Title)
	formatter.Println("URL:   %s", bookmark.URL)
	if len(bookmark.TagNames) > 0 {
		formatter.Println("Tags:  %s", strings.Join(bookmark.TagNames, ", "))
	}

	return nil
}
