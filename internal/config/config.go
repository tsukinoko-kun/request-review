package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var Version = 1

type Config struct {
	Version        int    `yaml:"version"`
	DiscordWebhook string `yaml:"discord_webhook,omitempty"`
}

func New() Config {
	return Config{
		Version: Version,
	}
}

func Load() (Config, error) {
	f, err := os.Open(".request-review.yaml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.Version > Version {
		return Config{}, fmt.Errorf("config version %d is not supported, update to use", cfg.Version)
	}

	return cfg, nil
}

func (cfg Config) Save() error {
	f, err := os.Create(".request-review.yaml")
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprint(f, "# yaml-language-server: https://raw.githubusercontent.com/tsukinoko-kun/request-review/main/src/config-schema.json\n---\n")
	return yaml.NewEncoder(f).Encode(cfg)
}
