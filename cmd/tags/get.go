package tags

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
	Use:   "get <id>",
	Short: "Get a tag by ID",
	Long:  "Retrieve detailed information about a specific tag.",
	Example: `  clinkding tags get 42`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func runGet(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid tag ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	tagsAPI := api.NewTagsAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	tag, err := tagsAPI.Get(ctx, id)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(tag)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(
			strconv.Itoa(tag.ID),
			tag.Name,
			strconv.Itoa(tag.BookmarkCount),
		)
		return nil
	}

	// Human-friendly output
	formatter.Println(formatter.Bold("Tag #%d"), tag.ID)
	formatter.Println("")
	formatter.Println("Name:       %s", tag.Name)
	formatter.Println("Bookmarks:  %d", tag.BookmarkCount)
	formatter.Println("Created:    %s", tag.DateAdded.Format(time.RFC3339))

	return nil
}
