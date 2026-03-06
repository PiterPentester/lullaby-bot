package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ckayt/lullaby/internal/bot"
	"github.com/ckayt/lullaby/internal/config"
	"github.com/ckayt/lullaby/internal/system"
	"github.com/joho/godotenv"
)

func main() {
	// Attempt to load .env file; it's okay if it doesn't exist (e.g., in production)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	hostRoot := os.Getenv("HOST_ROOT")
	if hostRoot == "" {
		hostRoot = "/host" // Default for k8s deployment
	}

	sys := system.NewManager(hostRoot)

	tgBot, err := bot.New(cfg, sys)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Channel to catch OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Bot started successfully")
		tgBot.Start()
	}()

	<-stop
	log.Println("Shutting down bot...")
	tgBot.Stop()
	log.Println("Stopped.")
}
