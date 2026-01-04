package assets

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:   "delete <bookmark-id> <asset-id>",
	Short: "Delete an asset",
	Long:  "Permanently delete an asset. This action cannot be undone.",
	Example: `  clinkding assets delete 42 1
  clinkding assets delete 42 1 --force`,
	Args: cobra.ExactArgs(2),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "skip confirmation")
}

func runDelete(cobraCmd *cobra.Command, args []string) error {
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

	// Get asset details for confirmation message
	asset, err := assetsAPI.Get(ctx, bookmarkID, assetID)
	if err != nil {
		return err
	}

	// Confirm deletion unless --force is used
	if !deleteForce && formatter.IsTTY() {
		fmt.Printf("Delete asset #%d \"%s\"? [y/N]: ", asset.ID, asset.DisplayName)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			formatter.Println("Aborted.")
			return nil
		}
	}

	// Delete asset
	if err := assetsAPI.Delete(ctx, bookmarkID, assetID); err != nil {
		return err
	}

	if !cfg.Quiet && !cfg.OutputJSON && !cfg.OutputPlain {
		formatter.Success("Asset #%d deleted", assetID)
	}

	return nil
}
