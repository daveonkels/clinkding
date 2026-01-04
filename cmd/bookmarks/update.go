package bookmarks

import (
	"context"
	"fmt"
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
	updateURL         string
	updateTitle       string
	updateDescription string
	updateNotes       string
	updateTags        string
	updateAddTags     string
	updateRemoveTags  string
	updateShared      string
	updateUnread      string
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a bookmark",
	Long:  "Update bookmark details. Only specified fields will be updated.",
	Example: `  clinkding bookmarks update 42 --title "New Title"
  clinkding bookmarks update 42 --add-tags "golang,tutorial"
  clinkding bookmarks update 42 --shared=true
  clinkding bookmarks update 42 --new-url "https://new-url.com" --title "Updated"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().StringVar(&updateURL, "new-url", "", "new URL for the bookmark")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "new title")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "new description")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "new notes")
	updateCmd.Flags().StringVar(&updateTags, "tags", "", "replace all tags (comma-separated)")
	updateCmd.Flags().StringVar(&updateAddTags, "add-tags", "", "add tags (comma-separated)")
	updateCmd.Flags().StringVar(&updateRemoveTags, "remove-tags", "", "remove tags (comma-separated)")
	updateCmd.Flags().StringVar(&updateShared, "shared", "", "set shared status (true/false)")
	updateCmd.Flags().StringVar(&updateUnread, "unread", "", "set unread status (true/false)")
}

func runUpdate(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()

	// Build update request
	update := &models.BookmarkUpdate{}

	if updateURL != "" {
		update.URL = updateURL
	}
	if updateTitle != "" {
		update.Title = updateTitle
	}
	if updateDescription != "" {
		update.Description = updateDescription
	}
	if updateNotes != "" {
		update.Notes = updateNotes
	}

	// Handle tag operations
	if updateTags != "" {
		tags := strings.Split(updateTags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		update.TagNames = tags
	} else if updateAddTags != "" || updateRemoveTags != "" {
		// Need to fetch current bookmark to manipulate tags
		current, err := bookmarksAPI.Get(ctx, id)
		if err != nil {
			return err
		}

		tagSet := make(map[string]bool)
		for _, tag := range current.TagNames {
			tagSet[tag] = true
		}

		if updateAddTags != "" {
			addTags := strings.Split(updateAddTags, ",")
			for _, tag := range addTags {
				tagSet[strings.TrimSpace(tag)] = true
			}
		}

		if updateRemoveTags != "" {
			removeTags := strings.Split(updateRemoveTags, ",")
			for _, tag := range removeTags {
				delete(tagSet, strings.TrimSpace(tag))
			}
		}

		var newTags []string
		for tag := range tagSet {
			newTags = append(newTags, tag)
		}
		update.TagNames = newTags
	}

	// Handle boolean flags
	if updateShared != "" {
		shared, err := strconv.ParseBool(updateShared)
		if err != nil {
			return fmt.Errorf("invalid --shared value: %s (use true or false)", updateShared)
		}
		update.Shared = &shared
	}
	if updateUnread != "" {
		unread, err := strconv.ParseBool(updateUnread)
		if err != nil {
			return fmt.Errorf("invalid --unread value: %s (use true or false)", updateUnread)
		}
		update.Unread = &unread
	}

	// Update bookmark
	bookmark, err := bookmarksAPI.Update(ctx, id, update)
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
	formatter.Success("Bookmark updated!")
	formatter.Println("")
	formatter.Println("ID:    %d", bookmark.ID)
	formatter.Println("Title: %s", bookmark.Title)
	formatter.Println("URL:   %s", bookmark.URL)
	if len(bookmark.TagNames) > 0 {
		formatter.Println("Tags:  %s", strings.Join(bookmark.TagNames, ", "))
	}

	return nil
}
