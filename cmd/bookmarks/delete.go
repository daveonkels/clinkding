package bookmarks

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
	Short: "Delete a bookmark",
	Long:  "Permanently delete a bookmark. This action cannot be undone.",
	Example: `  clinkding bookmarks delete 42
  clinkding bookmarks delete 42 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "skip confirmation")
}

func runDelete(cobraCmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %s", args[0])
	}

	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	bookmarksAPI := api.NewBookmarksAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()

	// Get bookmark details for confirmation message
	bookmark, err := bookmarksAPI.Get(ctx, id)
	if err != nil {
		return err
	}

	// Confirm deletion unless --force is used
	if !deleteForce && formatter.IsTTY() {
		fmt.Printf("Delete bookmark #%d \"%s\"? [y/N]: ", bookmark.ID, bookmark.Title)
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

	// Delete bookmark
	if err := bookmarksAPI.Delete(ctx, id); err != nil {
		return err
	}

	if !cfg.Quiet && !cfg.OutputJSON && !cfg.OutputPlain {
		formatter.Success("Bookmark #%d deleted", id)
	}

	return nil
}
