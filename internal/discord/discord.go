package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrWebhookNotSet = fmt.Errorf("webhook URL is not set")
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

func StartThread(webhook string, title string, content string) error {
	if webhook == "" {
		return ErrWebhookNotSet
	}
	wp := discordgo.WebhookParams{
		Content:    content,
		ThreadName: title,
	}
	sb := strings.Builder{}
	je := json.NewEncoder(&sb)
	je.Encode(wp)
	resp, err := http.Post(webhook, "application/json", strings.NewReader(sb.String()))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create thread: %s", string(body))
	}
	return nil
}
