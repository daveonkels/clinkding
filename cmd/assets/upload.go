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

var uploadCmd = &cobra.Command{
	Use:   "upload <bookmark-id> <file-path>",
	Short: "Upload an asset to a bookmark",
	Long:  "Upload a file as an asset attachment to a bookmark.",
	Example: `  clinkding assets upload 42 ~/Documents/screenshot.png
  clinkding assets upload 42 ./report.pdf`,
	Args: cobra.ExactArgs(2),
	RunE: runUpload,
}

func runUpload(cobraCmd *cobra.Command, args []string) error {
	bookmarkID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	filePath := args[1]

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	assetsAPI := api.NewAssetsAPI(httpClient)
	formatter := output.New(cfg)

	if !cfg.Quiet {
		formatter.Println("Uploading %s...", filePath)
	}

	ctx := context.Background()
	asset, err := assetsAPI.Upload(ctx, bookmarkID, filePath)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(asset)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(strconv.Itoa(asset.ID), asset.DisplayName)
		return nil
	}

	// Human-friendly output
	formatter.Success("Asset uploaded!")
	formatter.Println("")
	formatter.Println("Asset ID:    %d", asset.ID)
	formatter.Println("Bookmark ID: %d", asset.BookmarkID)
	formatter.Println("Name:        %s", asset.DisplayName)
	formatter.Println("Size:        %d KB", asset.FileSize/1024)

	return nil
}
