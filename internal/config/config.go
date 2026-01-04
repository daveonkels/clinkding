package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	URL         string
	Token       string
	OutputJSON  bool
	OutputPlain bool
	NoColor     bool
	Quiet       bool
	Verbose     bool

	// Default settings from config file
	Defaults struct {
		BookmarkLimit int    `mapstructure:"bookmark_limit"`
		OutputFormat  string `mapstructure:"output_format"`
	}
}

func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("defaults.bookmark_limit", 100)
	v.SetDefault("defaults.output_format", "auto")

	// Environment variables
	v.SetEnvPrefix("LINKDING")
	v.AutomaticEnv()

	// Config file
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		configDir, err := getConfigDir()
		if err != nil {
			return nil, err
		}
		v.AddConfigPath(configDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is okay, we can use env vars and flags
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Environment variables take precedence over config file
	if envURL := v.GetString("url"); envURL != "" {
		cfg.URL = envURL
	}
	if envToken := v.GetString("token"); envToken != "" {
		cfg.Token = envToken
	}

	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		cfg.NoColor = true
	}

	return cfg, nil
}

func getConfigDir() (string, error) {
	// Check LINKDING_CONFIG env var first
	if configPath := os.Getenv("LINKDING_CONFIG"); configPath != "" {
		return filepath.Dir(configPath), nil
	}

	// Use XDG_CONFIG_HOME or default to ~/.config
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(home, ".config")
	}

	return filepath.Join(configHome, "clinkding"), nil
}

func GetConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

func EnsureConfigDir() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}
