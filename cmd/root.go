package cmd

import (
	"fmt"

	"github.com/daveonkels/clinkding/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	url         string
	token       string
	outputJSON  bool
	outputPlain bool
	noColor     bool
	quiet       bool
	verbose     bool

	cfg         *config.Config
	version     string
	commit      string
	date        string
)

var rootCmd = &cobra.Command{
	Use:   "clinkding",
	Short: "Modern CLI for linkding bookmark manager",
	Long: `clinkding is a command-line interface for interacting with linkding,
a self-hosted bookmark management application.

Features:
  - Full API coverage: bookmarks, tags, bundles, assets, user profile
  - Human-friendly output with colors and tables
  - Script-friendly with --json and --plain flags
  - Safe operations with confirmations for destructive actions`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for certain commands
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override config with flags if provided
		if url != "" {
			cfg.URL = url
		}
		if token != "" {
			cfg.Token = token
		}

		// Set output preferences
		cfg.OutputJSON = outputJSON
		cfg.OutputPlain = outputPlain
		cfg.NoColor = noColor
		cfg.Quiet = quiet
		cfg.Verbose = verbose

		// Validate required config (except for config commands)
		if cmd.Parent() != nil && cmd.Parent().Name() != "config" {
			if cfg.URL == "" {
				return fmt.Errorf("linkding URL not configured. Use --url flag or run: clinkding config init")
			}
			if cfg.Token == "" {
				return fmt.Errorf("API token not configured. Use --token flag or run: clinkding config init")
			}
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: ~/.config/clinkding/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "linkding instance URL")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "API token")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVar(&outputPlain, "plain", false, "output as plain text")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colors")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.Version = version

	// Cobra's built-in completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

func RegisterCommands() {
	// Import subcommand packages here to avoid import cycles
	// They will be added via AddCommand
}

func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func GetConfig() *config.Config {
	return cfg
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
