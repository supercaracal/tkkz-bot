package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/supercaracal/tkkz-bot/internal/chats"
	"github.com/supercaracal/tkkz-bot/internal/configs"
	"github.com/supercaracal/tkkz-bot/internal/handlers"
	"github.com/supercaracal/tkkz-bot/internal/shared"
)

func waitUntil(
	chat chats.ChatClient,
	fail <-chan struct{},
	sign <-chan os.Signal,
) error {

	select {
	case <-sign:
		return chat.Disconnect()
	case <-fail:
		return fmt.Errorf("Failed to manage connection")
	}
}

func main() {
	godotenv.Load()

	logger := configs.NewBotLogger()
	cfg, err := configs.NewBotConfig()
	if err != nil {
		logger.Err.Fatalln(err)
	}
	ctx := shared.NewBotContext(cfg, logger)

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&ctx.Config.Verbose, "verbose", false, "verbose log")
	f.Parse(os.Args[1:])

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

	chat := chats.NewSlackClient(
		ctx.Config.SlackToken,
		ctx.Config.Verbose,
		ctx.Logger.Info,
	)

	h := handlers.NewEventHandler(ctx)
	chat.RegisterHandler("onMessage", h.RespondToContact)
	chat.RegisterHandler("onMessagingError", h.LogAsErr)
	chat.RegisterHandler("onAuthenticationError", h.LogAsErr)
	chat.RegisterHandler("onConnected", h.LogAsInfo)
	chat.RegisterHandler("onDisconnected", h.LogAsInfo)

	fail := make(chan struct{})
	chat.ConnectAsync(fail)
	chat.HandleEventsAsync()

	err = waitUntil(chat, fail, sign)
	if err != nil {
		ctx.Logger.Err.Fatalln(err)
	}

	ctx.Logger.Info.Println("Shutting Down")
	os.Exit(0)
}
