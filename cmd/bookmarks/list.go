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

var (
	listQuery         string
	listLimit         int
	listOffset        int
	listArchived      bool
	listUnarchived    bool
	listModifiedSince string
	listAddedSince    string
	listBundle        int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bookmarks",
	Long:  "List bookmarks with optional filters.",
	Example: `  clinkding bookmarks list
  clinkding bookmarks list --query "golang"
  clinkding bookmarks list --archived
  clinkding bookmarks list --modified-since "7d"
  clinkding bookmarks list --limit 20 --offset 40`,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVar(&listQuery, "query", "", "search query")
	listCmd.Flags().IntVar(&listLimit, "limit", 100, "max results")
	listCmd.Flags().IntVar(&listOffset, "offset", 0, "skip n results")
	listCmd.Flags().BoolVar(&listArchived, "archived", false, "show archived only")
	listCmd.Flags().BoolVar(&listUnarchived, "unarchived", false, "show unarchived only")
	listCmd.Flags().StringVar(&listModifiedSince, "modified-since", "", "filter by modification date (RFC3339 or relative: 24h, 7d)")
	listCmd.Flags().StringVar(&listAddedSince, "added-since", "", "filter by creation date (RFC3339 or relative: 24h, 7d)")
	listCmd.Flags().IntVar(&listBundle, "bundle", 0, "filter by bundle ID")
}

func runList(cobraCmd *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	// Parse date filters
	modifiedSince, err := parseDate(listModifiedSince)
	if err != nil {
		return fmt.Errorf("invalid --modified-since: %w", err)
	}
	addedSince, err := parseDate(listAddedSince)
	if err != nil {
		return fmt.Errorf("invalid --added-since: %w", err)
	}

	// Prepare list options
	opts := &api.ListOptions{
		Query:         listQuery,
		Limit:         listLimit,
		Offset:        listOffset,
		Archived:      listArchived,
		ModifiedSince: modifiedSince,
		AddedSince:    addedSince,
		BundleID:      listBundle,
	}

	ctx := context.Background()
	result, err := bookmarksAPI.List(ctx, opts)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(result)
	}

	if cfg.OutputPlain {
		for _, bookmark := range result.Results {
			output.PrintPlainLine(
				strconv.Itoa(bookmark.ID),
				bookmark.URL,
				bookmark.Title,
				strings.Join(bookmark.TagNames, ","),
			)
		}
		return nil
	}

	// Human-friendly table output
	if len(result.Results) == 0 {
		formatter.Info("No bookmarks found")
		return nil
	}

	table := output.NewTable([]string{"ID", "Title", "URL", "Tags", "Modified"})
	for _, bookmark := range result.Results {
		table.Append([]string{
			strconv.Itoa(bookmark.ID),
			output.TruncateString(bookmark.Title, 40),
			output.TruncateString(bookmark.URL, 50),
			output.FormatTags(bookmark.TagNames, 30),
			bookmark.DateModified.Format("2006-01-02"),
		})
	}
	table.Render()

	formatter.Println("")
	formatter.Println("Total: %d bookmarks", result.Count)
	if result.Next != nil {
		formatter.Info("Use --offset %d to see more", listOffset+listLimit)
	}

	return nil
}

func parseDate(dateStr string) (string, error) {
	if dateStr == "" {
		return "", nil
	}

	// Try parsing as RFC3339 first
	if _, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return dateStr, nil
	}

	// Try parsing as relative duration (24h, 7d, 30d, 1y)
	if strings.HasSuffix(dateStr, "h") || strings.HasSuffix(dateStr, "d") || strings.HasSuffix(dateStr, "y") {
		duration, err := parseRelativeDuration(dateStr)
		if err != nil {
			return "", err
		}
		past := time.Now().Add(-duration)
		return past.Format(time.RFC3339), nil
	}

	return "", fmt.Errorf("invalid date format (use RFC3339 or relative: 24h, 7d, 30d, 1y)")
}

func parseRelativeDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %w", err)
	}

	switch unit {
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "y":
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %s (use h, d, or y)", unit)
	}
}
