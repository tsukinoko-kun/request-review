package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/goccy/go-yaml"
	"github.com/tsukinoko-kun/request-review/internal/crypt"
	"github.com/tsukinoko-kun/request-review/internal/discord"
	"github.com/tsukinoko-kun/request-review/internal/git"
)

var Version = 1

type Config struct {
	Version              int    `yaml:"version"`
	name                 string `yaml:"-"`
	DiscordWebhook       string `yaml:"discord_webhook,omitempty"`
	LinearPersonalApiKey string `yaml:"linear_personal_api_key,omitempty"`
}

func New() Config {
	return Config{
		Version: Version,
		name:    git.RepoUrl(),
	}
}

func (cfg *Config) Edit() error {
	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Discord Webhook URL").
			Value(&cfg.DiscordWebhook).
			Validate(discord.ValidateWebhookURL).
			Placeholder("https://discord.com/api/webhooks/..."),
		huh.NewInput().
			Title("Linear Personal API Key").
			Description("Go to Settings > Security & Access > Personal API keys > New API key\nSelect Read permission").
			Value(&cfg.LinearPersonalApiKey).
			Placeholder("lin_api_..."),
	)).Run(); err != nil {
		return err
	}
	return nil
}

func Load() (Config, error) {
	f, err := os.Open(".request-review.yaml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	cfg := Config{
		name: git.RepoUrl(),
	}
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.Version > Version {
		return Config{}, fmt.Errorf("config version %d is not supported, update to use", cfg.Version)
	}

	cfg.DiscordWebhook, err = crypt.Decrypt(cfg.name, cfg.DiscordWebhook)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decrypt Discord webhook URL: %w", err)
	}
	cfg.LinearPersonalApiKey, err = crypt.Decrypt(cfg.name, cfg.LinearPersonalApiKey)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decrypt Linear personal API key: %w", err)
	}

	return cfg, nil
}

func (cfg Config) Save() error {
	var err error
	cfg.DiscordWebhook, err = crypt.Encrypt(cfg.name, cfg.DiscordWebhook)
	if err != nil {
		return fmt.Errorf("failed to encrypt Discord webhook URL: %w", err)
	}
	cfg.LinearPersonalApiKey, err = crypt.Encrypt(cfg.name, cfg.LinearPersonalApiKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt Linear personal API key: %w", err)
	}

	f, err := os.Create(".request-review.yaml")
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprint(f, "# yaml-language-server: https://raw.githubusercontent.com/tsukinoko-kun/request-review/main/src/config-schema.json\n---\n")
	return yaml.NewEncoder(f).Encode(cfg)
}
