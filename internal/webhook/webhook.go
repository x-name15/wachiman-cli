package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Send(webhookURL, webhookType, title, message string) error {
	if webhookURL == "" {
		return nil
	}

	var payload interface{}

	switch webhookType {
	case "discord":
		payload = buildDiscord(title, message)
	default: 
		payload = buildSlack(title, message)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error mandando webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook respondió con %d", resp.StatusCode)
	}

	return nil
}

func buildSlack(title, message string) map[string]interface{} {
	return map[string]interface{}{
		"username": "wachiman",
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s*\n%s", title, message),
				},
			},
			{
				"type": "context",
				"elements": []map[string]interface{}{
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("wachiman · %s", time.Now().Format("2006-01-02 15:04:05")),
					},
				},
			},
		},
	}
}

func buildDiscord(title, message string) map[string]interface{} {
	color := 3066993 // verde
	if containsAny(title, []string{"CAÍDO", "caído", "error", "Error"}) {
		color = 15158332 
	}

	return map[string]interface{}{
		"username": "wachiman",
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"description": message,
				"color":       color,
				"footer": map[string]interface{}{
					"text": fmt.Sprintf("wachiman · %s", time.Now().Format("2006-01-02 15:04:05")),
				},
			},
		},
	}
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}