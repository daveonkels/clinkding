package bundles

import (
	"context"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bundles",
	Long:  "List all bundles.",
	Example: `  clinkding bundles list
  clinkding bundles list --json`,
	RunE: runList,
}

func runList(cobraCmd *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bundlesAPI := api.NewBundlesAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	result, err := bundlesAPI.List(ctx)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(result)
	}

	if cfg.OutputPlain {
		for _, bundle := range result.Results {
			output.PrintPlainLine(
				strconv.Itoa(bundle.ID),
				bundle.Name,
				bundle.Description,
			)
		}
		return nil
	}

	// Human-friendly table output
	if len(result.Results) == 0 {
		formatter.Info("No bundles found")
		return nil
	}

	table := output.NewTable([]string{"ID", "Name", "Description", "Created"})
	for _, bundle := range result.Results {
		table.Append([]string{
			strconv.Itoa(bundle.ID),
			bundle.Name,
			output.TruncateString(bundle.Description, 50),
			bundle.DateAdded.Format("2006-01-02"),
		})
	}
	table.Render()

	formatter.Println("")
	formatter.Println("Total: %d bundles", result.Count)

	return nil
}
