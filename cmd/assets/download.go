package assets

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var downloadOutput string

var downloadCmd = &cobra.Command{
	Use:   "download <bookmark-id> <asset-id> [output-path]",
	Short: "Download an asset",
	Long:  "Download an asset file from a bookmark.",
	Example: `  clinkding assets download 42 1
  clinkding assets download 42 1 ~/Downloads/myfile.pdf
  clinkding assets download 42 1 -o ./screenshot.png`,
	Args: cobra.RangeArgs(2, 3),
	RunE: runDownload,
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadOutput, "output", "o", "", "output file path")
}

func runDownload(cobraCmd *cobra.Command, args []string) error {
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

	// Determine output path
	outputPath := downloadOutput
	if outputPath == "" && len(args) >= 3 {
		outputPath = args[2]
	}

	// If still no output path, get asset details to use its display name
	if outputPath == "" {
		asset, err := assetsAPI.Get(ctx, bookmarkID, assetID)
		if err != nil {
			return err
		}
		outputPath = asset.DisplayName
	}

	// Make sure we have an absolute or relative path
	if outputPath == "" {
		outputPath = fmt.Sprintf("asset-%d", assetID)
	}

	if !cfg.Quiet {
		formatter.Println("Downloading to %s...", outputPath)
	}

	// Download the file
	if err := assetsAPI.Download(ctx, bookmarkID, assetID, outputPath); err != nil {
		return err
	}

	if !cfg.Quiet && !cfg.OutputJSON && !cfg.OutputPlain {
		absPath, _ := filepath.Abs(outputPath)
		formatter.Success("Asset downloaded to: %s", absPath)
	}

	return nil
}
