package shared

import (
	"github.com/supercaracal/tkkz-bot/internal/config"
)

// BotContext is
type BotContext struct {
	Config *config.BotConfig
	Logger *config.BotLogger
}

// NewBotContext is
func NewBotContext(c *config.BotConfig, l *config.BotLogger) *BotContext {
	return &BotContext{Config: c, Logger: l}
}
