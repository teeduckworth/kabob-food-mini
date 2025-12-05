package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/caarlos0/env/v10"

	"github.com/rashidmailru/kabobfood/internal/bot"
)

type config struct {
	TelegramToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	MiniAppURL    string `env:"MINI_APP_URL" envDefault:"https://kabob-food-mini.vercel.app"`
	Debug         bool   `env:"BOT_DEBUG" envDefault:"false"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	backendURL := strings.TrimSpace(os.Getenv("BOT_BACKEND_URL"))
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}
	log.Printf("bot backend url: %s", backendURL)

	service, err := bot.New(bot.Config{
		Token:      cfg.TelegramToken,
		BackendURL: backendURL,
		MiniAppURL: cfg.MiniAppURL,
		Debug:      cfg.Debug,
	})
	if err != nil {
		log.Fatalf("failed to init bot: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := service.Run(ctx); err != nil {
		log.Fatalf("bot stopped with error: %v", err)
	}
}
