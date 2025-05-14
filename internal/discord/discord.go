package discord

import (
	"fmt"
	"strings"
)

func ValidateWebhookURL(url string) error {
	if url == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}
	if !strings.HasPrefix(url, "https://discord.com/api/webhooks/") {
		return fmt.Errorf("invalid webhook URL: %s", url)
	}
	return nil
}
