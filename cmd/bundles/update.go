package bundles

import (
	"context"
	"fmt"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var (
	updateName        string
	updateDescription string
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a bundle",
	Long:  "Update bundle name and/or description. Only specified fields will be updated.",
	Example: `  clinkding bundles update 42 --name "New Name"
  clinkding bundles update 42 --description "Updated description"
  clinkding bundles update 42 --name "Go Lang" --description "Everything Go"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().StringVar(&updateName, "name", "", "new bundle name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "new description")
}

func runUpdate(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bundle ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bundlesAPI := api.NewBundlesAPI(httpClient)
	formatter := output.New(cfg)

	update := &models.BundleUpdate{
		Name:        updateName,
		Description: updateDescription,
	}

	ctx := context.Background()
	bundle, err := bundlesAPI.Update(ctx, id, update)
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
	formatter.Success("Bundle updated!")
	formatter.Println("")
	formatter.Println("ID:   %d", bundle.ID)
	formatter.Println("Name: %s", bundle.Name)
	if bundle.Description != "" {
		formatter.Println("Description: %s", bundle.Description)
	}

	return nil
}
