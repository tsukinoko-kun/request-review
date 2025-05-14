package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/tsukinoko-kun/request-review/internal/config"
	"github.com/tsukinoko-kun/request-review/internal/discord"
	"github.com/tsukinoko-kun/request-review/internal/git"
)

func main() {
	if len(os.Args) != 2 {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s <from>..<to>\n", exe)
		os.Exit(1)
		return
	}

	if os.Args[1] == "init" {
		initCfg()
	} else {
		requestReview()
	}
}

func initCfg() {
	cfg := config.New()
	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Discord Webhook URL").
			Value(&cfg.DiscordWebhook).
			Validate(discord.ValidateWebhookURL).
			Placeholder("https://discord.com/api/webhooks/..."),
	)).Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
			return
		}
	}
	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
}

func requestReview() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
	splitFromTo := strings.SplitN(os.Args[1], "..", 2)
	if patch, err := git.Patch(cfg, splitFromTo[0], splitFromTo[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	} else {
		fmt.Println(patch)
	}
}
