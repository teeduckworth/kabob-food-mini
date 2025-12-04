package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v10"

	"github.com/rashidmailru/kabobfood/internal/bot"
)

type config struct {
	TelegramToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	BackendURL    string `env:"BOT_BACKEND_URL" envDefault:"http://localhost:8080"`
	MiniAppURL    string `env:"MINI_APP_URL" envDefault:"https://kabob-food-mini.vercel.app"`
	Debug         bool   `env:"BOT_DEBUG" envDefault:"false"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	service, err := bot.New(bot.Config{
		Token:      cfg.TelegramToken,
		BackendURL: cfg.BackendURL,
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
