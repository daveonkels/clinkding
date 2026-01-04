package tags

import (
	"context"
	"strconv"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new tag",
	Long:  "Create a new tag with the specified name.",
	Example: `  clinkding tags create golang
  clinkding tags create "machine learning"`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func runCreate(cobraCmd *cobra.Command, args []string) error {
	tagName := args[0]

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	tagsAPI := api.NewTagsAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	tag, err := tagsAPI.Create(ctx, tagName)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(tag)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(strconv.Itoa(tag.ID), tag.Name)
		return nil
	}

	// Human-friendly output
	formatter.Success("Tag created!")
	formatter.Println("")
	formatter.Println("ID:   %d", tag.ID)
	formatter.Println("Name: %s", tag.Name)

	return nil
}
