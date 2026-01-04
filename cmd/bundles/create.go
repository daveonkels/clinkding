package bundles

import (
	"context"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var (
	createDescription string
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new bundle",
	Long:  "Create a new bundle with the specified name and optional description.",
	Example: `  clinkding bundles create "Go Resources"
  clinkding bundles create "AI & ML" --description "AI and machine learning bookmarks"`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVar(&createDescription, "description", "", "bundle description")
}

func runCreate(cobraCmd *cobra.Command, args []string) error {
	bundleName := args[0]

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bundlesAPI := api.NewBundlesAPI(httpClient)
	formatter := output.New(cfg)

	bundleCreate := &models.BundleCreate{
		Name:        bundleName,
		Description: createDescription,
	}

	ctx := context.Background()
	bundle, err := bundlesAPI.Create(ctx, bundleCreate)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(bundle)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(strconv.Itoa(bundle.ID), bundle.Name)
		return nil
	}

	// Human-friendly output
	formatter.Success("Bundle created!")
	formatter.Println("")
	formatter.Println("ID:   %d", bundle.ID)
	formatter.Println("Name: %s", bundle.Name)
	if bundle.Description != "" {
		formatter.Println("Description: %s", bundle.Description)
	}

	return nil
}
