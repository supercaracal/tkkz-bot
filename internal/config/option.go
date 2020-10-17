package config

import (
	"fmt"
	"os"
)

// BotOption is
type BotOption struct {
	SlackToken string
	BotID      string
	BrainURL   string
	Verbose    bool
	Debug      bool
}

// NewBotOption is
func NewBotOption() (*BotOption, error) {
	opt := BotOption{
		SlackToken: os.Getenv("SLACK_TOKEN"),
		BotID:      os.Getenv("BOT_ID"),
		BrainURL:   os.Getenv("BRAIN_URL"),
	}

	if opt.SlackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN environment variable required")
	}

	if opt.BotID == "" {
		return nil, fmt.Errorf("BOT_ID environment variable required")
	}

	return &opt, nil
}
