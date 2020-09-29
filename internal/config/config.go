package config

import (
	"fmt"
	"os"
)

// BotConfig is
type BotConfig struct {
	SlackToken string
	BotID      string
	Verbose    bool
}

// NewBotConfig is
func NewBotConfig() (*BotConfig, error) {
	cfg := BotConfig{
		SlackToken: os.Getenv("SLACK_TOKEN"),
		BotID:      os.Getenv("BOT_ID"),
	}

	if cfg.SlackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN environment variable required")
	}

	if cfg.BotID == "" {
		return nil, fmt.Errorf("BOT_ID environment variable required")
	}

	return &cfg, nil
}
