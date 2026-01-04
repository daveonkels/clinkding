package assets

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get <bookmark-id> <asset-id>",
	Short:   "Get asset details",
	Long:    "Retrieve detailed information about a specific asset.",
	Example: `  clinkding assets get 42 1`,
	Args:    cobra.ExactArgs(2),
	RunE:    runGet,
}

func runGet(cobraCmd *cobra.Command, args []string) error {
	bookmarkID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	assetID, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid asset ID: %s", args[1])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	assetsAPI := api.NewAssetsAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	asset, err := assetsAPI.Get(ctx, bookmarkID, assetID)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(asset)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(
			strconv.Itoa(asset.ID),
			asset.DisplayName,
			strconv.FormatInt(asset.FileSize, 10),
		)
		return nil
	}

	// Human-friendly output
	formatter.Println(formatter.Bold("Asset #%d"), asset.ID)
	formatter.Println("")
	formatter.Println("Bookmark ID: %d", asset.BookmarkID)
	formatter.Println("Name:        %s", asset.DisplayName)
	formatter.Println("File:        %s", asset.File)
	formatter.Println("Size:        %d KB", asset.FileSize/1024)
	formatter.Println("Status:      %s", asset.Status)
	formatter.Println("Created:     %s", asset.DateCreated.Format(time.RFC3339))

	return nil
}
