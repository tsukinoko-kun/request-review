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
	"github.com/tsukinoko-kun/request-review/internal/linear"
	"github.com/tsukinoko-kun/request-review/internal/metadata"
)

func main() {
	switch len(os.Args) {
	case 1:
		requestReviewSmart()
	case 2:
		switch os.Args[1] {
		case "help":
			help()
		case "init":
			initCfg()
		case "smart":
			requestReviewSmart()
		case "version":
			fmt.Println(metadata.Version)
		default:
			if strings.Contains(os.Args[1], "..") {
				splitFromTo := strings.SplitN(os.Args[1], "..", 2)
				requestReviewRange(splitFromTo[0], splitFromTo[1])
			}
			help()
			os.Exit(1)
		}
	default:
		help()
		os.Exit(1)
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
		huh.NewInput().
			Title("Linear Personal API Key").
			Description("Go to Settings > Security & Access > Personal API keys > New API key\nSelect Read permission").
			Value(&cfg.LinearPersonalApiKey).
			Placeholder("lin_api_..."),
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

func requestReviewSmart() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
	if patch, err := git.SmartPatch(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	} else {
		requestReviewForPatch(cfg, patch, "")
	}
}

func requestReviewRange(from, to string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
	if patch, err := git.Patch(cfg, from, to); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	} else {
		requestReviewForPatch(cfg, patch, os.Args[1])
	}
}

func requestReviewForPatch(cfg config.Config, patch string, label string) {
	fi, err := git.GetRepoInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
	var (
		title string
		body  string
	)
	author := git.User()
	if issue, err := linear.FindIssueByBranchName(cfg, fi.Bookmark()); err == nil {
		title = fmt.Sprintf("%s in %s @ %s", issue.Title, fi.Bookmark(), fi.Name())
		body = fmt.Sprintf("%s requested review for %s\n\n%s\n\n```diff\n%s\n```", author, issue.Title, issue.Description, patch)
	} else if label == "" {
		title = fmt.Sprintf("%s @ %s", fi.Bookmark(), fi.Name())
		body = fmt.Sprintf("%s request review for %s\n```diff\n%s\n```", author, fi.Bookmark(), patch)
	} else {
		title = fmt.Sprintf("%s in %s @ %s", label, fi.Bookmark(), fi.Name())
		body = fmt.Sprintf("%s request review for %s\n```diff\n%s\n```", author, label, patch)
	}

	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Title").Value(&title).Validate(huh.ValidateNotEmpty()),
		huh.NewText().Title("Body").Value(&body).Validate(huh.ValidateNotEmpty()),
	)).Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}

	if err := discord.StartThread(cfg, title, body); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
}

func help() {
	exe := filepath.Base(os.Args[0])
	fmt.Printf("%s <from>..<to> - Request review for a range of commits", exe)
	fmt.Printf("%s smart        - Request review with smart range", exe)
	fmt.Printf("%s init         - Initialize a configuration", exe)
	fmt.Printf("%s help         - Show this help message", exe)
}
