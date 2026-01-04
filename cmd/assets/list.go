package assets

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

var listCmd = &cobra.Command{
	Use:   "list <bookmark-id>",
	Short: "List assets for a bookmark",
	Long:  "List all file assets attached to a specific bookmark.",
	Example: `  clinkding assets list 42
  clinkding assets list 42 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runList,
}

func runList(cobraCmd *cobra.Command, args []string) error {
	bookmarkID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	assetsAPI := api.NewAssetsAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	result, err := assetsAPI.List(ctx, bookmarkID)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(result)
	}

	if cfg.OutputPlain {
		for _, asset := range result.Results {
			output.PrintPlainLine(
				strconv.Itoa(asset.ID),
				asset.DisplayName,
				strconv.FormatInt(asset.FileSize, 10),
				asset.Status,
			)
		}
		return nil
	}

	// Human-friendly table output
	if len(result.Results) == 0 {
		formatter.Info("No assets found for bookmark #%d", bookmarkID)
		return nil
	}

	table := output.NewTable([]string{"ID", "Name", "Size", "Status", "Created"})
	for _, asset := range result.Results {
		sizeKB := asset.FileSize / 1024
		table.Append([]string{
			strconv.Itoa(asset.ID),
			asset.DisplayName,
			fmt.Sprintf("%d KB", sizeKB),
			asset.Status,
			asset.DateCreated.Format("2006-01-02"),
		})
	}
	table.Render()

	formatter.Println("")
	formatter.Println("Total: %d assets", result.Count)

	return nil
}
