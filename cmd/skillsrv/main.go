package main

import (
	"AliceSkills/internal/app"
	"AliceSkills/pkg/config"
	"AliceSkills/pkg/openai"
	"context"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	openai.Init(openai.Config{
		Model:      "",
		Timeout:    30 * time.Second,
		MaxRetries: 2,
		System:     string(openai.VoiceHelperSystemPrompt),
	})

	a, err := app.New(cfg)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}
