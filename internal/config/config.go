package config

import (
	"fmt"
	"os"
)

// BotConfig is
type BotConfig struct {
	SlackToken string
	BotID      string
	BrainURL   string
	Verbose    bool
	Debug      bool
}

// NewBotConfig is
func NewBotConfig() (*BotConfig, error) {
	cfg := BotConfig{
		SlackToken: os.Getenv("SLACK_TOKEN"),
		BotID:      os.Getenv("BOT_ID"),
		BrainURL:   os.Getenv("BRAIN_URL"),
	}

	if cfg.SlackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN environment variable required")
	}

	if cfg.BotID == "" {
		return nil, fmt.Errorf("BOT_ID environment variable required")
	}

	return &cfg, nil
}
