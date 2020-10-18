package config

// BotContext is
type BotContext struct {
	Option *BotOption
	Logger *BotLogger
}

// NewBotContext is
func NewBotContext(o *BotOption, l *BotLogger) *BotContext {
	return &BotContext{Option: o, Logger: l}
}
