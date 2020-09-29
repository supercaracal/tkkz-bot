package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/supercaracal/tkkz-bot/internal/chat"
	"github.com/supercaracal/tkkz-bot/internal/config"
	"github.com/supercaracal/tkkz-bot/internal/handler"
	"github.com/supercaracal/tkkz-bot/internal/shared"
)

func waitUntil(
	chatCli chat.Client,
	fail <-chan struct{},
	sign <-chan os.Signal,
) error {

	select {
	case <-sign:
		return chatCli.Disconnect()
	case <-fail:
		return fmt.Errorf("Failed to manage connection")
	}
}

func main() {
	godotenv.Load()

	logger := config.NewBotLogger()
	cfg, err := config.NewBotConfig()
	if err != nil {
		logger.Err.Fatalln(err)
	}
	ctx := shared.NewBotContext(cfg, logger)

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&ctx.Config.Verbose, "verbose", false, "verbose log")
	f.Parse(os.Args[1:])

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

	chatCli := chat.NewSlackClient(
		ctx.Config.SlackToken,
		ctx.Config.Verbose,
		ctx.Logger.Info,
	)

	h := handler.NewEventHandler(ctx)
	chatCli.RegisterHandler("onMessage", h.RespondToContact)
	chatCli.RegisterHandler("onMessagingError", h.LogAsErr)
	chatCli.RegisterHandler("onAuthenticationError", h.LogAsErr)
	chatCli.RegisterHandler("onConnected", h.LogAsInfo)
	chatCli.RegisterHandler("onDisconnected", h.LogAsInfo)

	fail := make(chan struct{})
	chatCli.ConnectAsync(fail)
	chatCli.HandleEventsAsync()

	err = waitUntil(chatCli, fail, sign)
	if err != nil {
		ctx.Logger.Err.Fatalln(err)
	}

	ctx.Logger.Info.Println("Shutting Down")
	os.Exit(0)
}
