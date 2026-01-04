package config

import (
	"fmt"
	"os"

	"github.com/daveonkels/clinkding/cmd"
	"github.com/daveonkels/clinkding/internal/config"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the merged configuration from all sources (file, environment, flags).",
	RunE:  runShow,
}

func runShow(cmdCobra *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	if cfg == nil {
		// Load config manually if not loaded
		var err error
		cfg, err = config.Load("")
		if err != nil {
			return err
		}
	}

	configPath, _ := config.GetConfigPath()

	fmt.Println("Configuration:")
	fmt.Printf("  Config file: %s\n", configPath)
	fmt.Printf("  URL: %s\n", cfg.URL)
	fmt.Printf("  Token: %s\n", redactToken(cfg.Token))
	fmt.Printf("  Default bookmark limit: %d\n", cfg.Defaults.BookmarkLimit)
	fmt.Printf("  Default output format: %s\n", cfg.Defaults.OutputFormat)
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  LINKDING_URL:", getEnvOrNone("LINKDING_URL"))
	fmt.Println("  LINKDING_TOKEN:", redactToken(getEnvOrNone("LINKDING_TOKEN")))
	fmt.Println("  NO_COLOR:", getEnvOrNone("NO_COLOR"))

	return nil
}

func redactToken(token string) string {
	if token == "" {
		return "(not set)"
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

func getEnvOrNone(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return "(not set)"
}
