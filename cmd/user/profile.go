package user

import (
	"context"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/output"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show user profile",
	Long:  "Display user profile and preferences.",
	Example: `  clinkding user profile
  clinkding user profile --json`,
	RunE: runProfile,
}

func runProfile(cobraCmd *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	httpClient := client.New(cfg.URL, cfg.Token)
	userAPI := api.NewUserAPI(httpClient)
	formatter := output.New(cfg)

	ctx := context.Background()
	profile, err := userAPI.GetProfile(ctx)
	if err != nil {
		return err
	}

	// Output based on format
	if cfg.OutputJSON {
		return formatter.PrintJSON(profile)
	}

	if cfg.OutputPlain {
		output.PrintPlainLine(
			profile.Theme,
			profile.BookmarkDateDisplay,
			boolToString(profile.EnableSharing),
			boolToString(profile.EnableFavicons),
		)
		return nil
	}

	// Human-friendly output
	formatter.Println(formatter.Bold("User Profile"))
	formatter.Println("")
	formatter.Println("Appearance:")
	formatter.Println("  Theme:                  %s", profile.Theme)
	formatter.Println("  Bookmark date display:  %s", profile.BookmarkDateDisplay)
	formatter.Println("  Bookmark link target:   %s", profile.BookmarkLinkTarget)
	formatter.Println("")
	formatter.Println("Features:")
	formatter.Println("  Sharing enabled:        %t", profile.EnableSharing)
	formatter.Println("  Public sharing:         %t", profile.EnablePublicSharing)
	formatter.Println("  Favicons enabled:       %t", profile.EnableFavicons)
	formatter.Println("  Preview images:         %t", profile.EnablePreviewImages)
	formatter.Println("  Display URL:            %t", profile.DisplayURL)
	formatter.Println("  Display viewed date:    %t", profile.DisplayViewedDate)
	formatter.Println("  Permanent notes:        %t", profile.PermanentNotes)
	formatter.Println("")
	formatter.Println("Web Archive:")
	formatter.Println("  Integration:            %s", profile.WebArchiveIntegration)
	formatter.Println("")
	formatter.Println("Search Preferences:")
	formatter.Println("  Sort:                   %s", profile.SearchPreferences.Sort)
	formatter.Println("  Shared bookmarks:       %s", profile.SearchPreferences.SharedBookmarks)
	formatter.Println("  Unread only:            %t", profile.SearchPreferences.UnreadOnly)

	return nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
