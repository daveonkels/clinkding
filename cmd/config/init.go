package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/daveonkels/clinkding/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  "Create a new configuration file with prompts for URL and API token.",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	// Ensure config directory exists
	_, err := config.EnsureConfigDir()
	if err != nil {
		return err
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at: %s\n", configPath)
		fmt.Print("Overwrite? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// Prompt for configuration
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Configure clinkding CLI")
	fmt.Println()

	fmt.Print("Linkding URL: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	url = strings.TrimSpace(url)

	fmt.Print("API Token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	token = strings.TrimSpace(token)

	// Create config structure
	cfg := map[string]interface{}{
		"url":   url,
		"token": token,
		"defaults": map[string]interface{}{
			"bookmark_limit": 100,
			"output_format":  "auto",
		},
	}

	// Write config file
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() { _ = file.Close() }()

	encoder := yaml.NewEncoder(file)
	defer func() { _ = encoder.Close() }()

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("\n✓ Configuration saved to: %s\n", configPath)
	fmt.Println("\nTest your configuration with:")
	fmt.Println("  clinkding config test")

	return nil
}
