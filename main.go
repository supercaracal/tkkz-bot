package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/supercaracal/tkkz-bot/internal/configurations"
	"github.com/supercaracal/tkkz-bot/internal/handlers"
	"github.com/supercaracal/tkkz-bot/internal/infrastructures"
	"github.com/supercaracal/tkkz-bot/internal/interfaces"
)

func waitUntilShutdown(chat interfaces.Chat, fail <-chan struct{}, sign <-chan os.Signal) error {
	select {
	case <-sign:
		return chat.Disconnect()
	case <-fail:
		return fmt.Errorf("Failed to manage connection")
	}
}

func main() {
	godotenv.Load()
	logger := configurations.NewMyLogger()

	cfg, err := configurations.NewMyConfig()
	if err != nil {
		logger.Err.Fatalln(err)
	}

	cfg.Logger = logger

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&cfg.Verbose, "verbose", false, "verbose log")
	f.Parse(os.Args[1:])

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

	var chat interfaces.Chat = infrastructures.NewSlackClient(cfg.SlackToken, cfg.Verbose, cfg.Logger.Info)

	h := handlers.NewEventHandler(cfg)
	chat.RegisterHandler("onMessage", h.RespondToContact)
	chat.RegisterHandler("onMessagingError", h.LogAsErr)
	chat.RegisterHandler("onAuthenticationError", h.LogAsErr)
	chat.RegisterHandler("onConnected", h.LogAsInfo)
	chat.RegisterHandler("onDisconnected", h.LogAsInfo)

	fail := make(chan struct{})
	chat.ConnectAsync(fail)
	chat.HandleEventsAsync()

	err = waitUntilShutdown(chat, fail, sign)
	if err != nil {
		cfg.Logger.Err.Fatalln(err)
	}

	cfg.Logger.Info.Println("Shutting Down")
	os.Exit(0)
}
