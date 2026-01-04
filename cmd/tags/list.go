package tags

import (
	"context"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listOffset int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
	Long:  "List all tags with pagination support.",
	Example: `  clinkding tags list
  clinkding tags list --limit 50 --offset 100`,
	RunE: runList,
}

func init() {
	listCmd.Flags().IntVar(&listLimit, "limit", 100, "max results")
	listCmd.Flags().IntVar(&listOffset, "offset", 0, "skip n results")
}

func runList(cobraCmd *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	tagsAPI := api.NewTagsAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	result, err := tagsAPI.List(ctx, listLimit, listOffset)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(result)
	}

	if cfg.OutputPlain {
		for _, tag := range result.Results {
			output.PrintPlainLine(
				strconv.Itoa(tag.ID),
				tag.Name,
				strconv.Itoa(tag.BookmarkCount),
			)
		}
		return nil
	}

	// Human-friendly table output
	if len(result.Results) == 0 {
		formatter.Info("No tags found")
		return nil
	}

	table := output.NewTable([]string{"ID", "Name", "Bookmarks", "Created"})
	for _, tag := range result.Results {
		table.Append([]string{
			strconv.Itoa(tag.ID),
			tag.Name,
			strconv.Itoa(tag.BookmarkCount),
			tag.DateAdded.Format("2006-01-02"),
		})
	}
	table.Render()

	formatter.Println("")
	formatter.Println("Total: %d tags", result.Count)
	if result.Next != nil {
		formatter.Info("Use --offset %d to see more", listOffset+listLimit)
	}

	return nil
}
