package config

import (
	"fmt"
	"os"
)

// BotOption is
type BotOption struct {
	SlackToken    string
	BotID         string
	SlackAppToken string
	SlackBotToken string
	BrainURL      string
	Verbose       bool
	Debug         bool
}

// NewBotOption is
func NewBotOption() (*BotOption, error) {
	opt := BotOption{
		SlackToken:    os.Getenv("SLACK_TOKEN"),
		BotID:         os.Getenv("BOT_ID"),
		SlackAppToken: os.Getenv("SLACK_APP_TOKEN"),
		SlackBotToken: os.Getenv("SLACK_BOT_TOKEN"),
		BrainURL:      os.Getenv("BRAIN_URL"),
	}

	if opt.SlackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN environment variable required")
	}

	if opt.BotID == "" {
		return nil, fmt.Errorf("BOT_ID environment variable required")
	}

	if opt.SlackAppToken == "" {
		return nil, fmt.Errorf("SLACK_APP_TOKEN environment variable required")
	}

	if opt.SlackBotToken == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN environment variable required")
	}

	return &opt, nil
}
