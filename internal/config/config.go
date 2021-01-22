package config

// BotConfig is
type BotConfig struct {
	Option *BotOption
	Logger *BotLogger
}

// NewBotConfig is
func NewBotConfig(o *BotOption, l *BotLogger) *BotConfig {
	return &BotConfig{Option: o, Logger: l}
}
