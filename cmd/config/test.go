package config

import (
	"context"
	"fmt"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/api"
	"github.com/daveonkels/clinkding/internal/client"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test connection to linkding",
	Long:  "Validate URL and API token by fetching user profile.",
	RunE:  runTest,
}

func runTest(cmdCobra *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if cfg.URL == "" {
		return fmt.Errorf("linkding URL not configured")
	}
	if cfg.Token == "" {
		return fmt.Errorf("API token not configured")
	}

	fmt.Printf("Testing connection to %s...\n", cfg.URL)

	// Create client and API
	httpClient := client.New(cfg.URL, cfg.Token)
	userAPI := api.NewUserAPI(httpClient)

	// Fetch user profile
	ctx := context.Background()
	profile, err := userAPI.GetProfile(ctx)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	fmt.Println("\n✓ Connection successful!")
	fmt.Println("\nUser profile:")
	fmt.Printf("  Theme: %s\n", profile.Theme)
	fmt.Printf("  Bookmark date display: %s\n", profile.BookmarkDateDisplay)
	fmt.Printf("  Sharing enabled: %t\n", profile.EnableSharing)
	fmt.Printf("  Favicons enabled: %t\n", profile.EnableFavicons)

	return nil
}
