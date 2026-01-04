package bundles

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
	Use:     "get <id>",
	Short:   "Get a bundle by ID",
	Long:    "Retrieve detailed information about a specific bundle.",
	Example: `  clinkding bundles get 42`,
	Args:    cobra.ExactArgs(1),
	RunE:    runGet,
}

func runGet(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bundle ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bundlesAPI := api.NewBundlesAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	bundle, err := bundlesAPI.Get(ctx, id)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(bundle)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(
			strconv.Itoa(bundle.ID),
			bundle.Name,
			bundle.Description,
		)
		return nil
	}

	// Human-friendly output
	formatter.Println(formatter.Bold("Bundle #%d"), bundle.ID)
	formatter.Println("")
	formatter.Println("Name:        %s", bundle.Name)
	if bundle.Description != "" {
		formatter.Println("Description: %s", bundle.Description)
	}
	formatter.Println("Created:     %s", bundle.DateAdded.Format(time.RFC3339))

	return nil
}
