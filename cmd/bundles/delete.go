package bundles

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
	Use:   "delete <id>",
	Short: "Delete a bundle",
	Long:  "Permanently delete a bundle. This action cannot be undone.",
	Example: `  clinkding bundles delete 42
  clinkding bundles delete 42 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "skip confirmation")
}

func runDelete(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bundle ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bundlesAPI := api.NewBundlesAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()

	// Get bundle details for confirmation message
	bundle, err := bundlesAPI.Get(ctx, id)
	if err != nil {
		return err
	}

	// Confirm deletion unless --force is used
	if !deleteForce && formatter.IsTTY() {
		fmt.Printf("Delete bundle #%d \"%s\"? [y/N]: ", bundle.ID, bundle.Name)
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

	// Delete bundle
	if err := bundlesAPI.Delete(ctx, id); err != nil {
		return err
	}

	if !cfg.Quiet && !cfg.OutputJSON && !cfg.OutputPlain {
		formatter.Success("Bundle #%d deleted", id)
	}

	return nil
}
