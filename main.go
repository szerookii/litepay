package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/szerookii/litepay/cryptocurrency"
	"github.com/szerookii/litepay/db"
	"github.com/szerookii/litepay/router"
	"os"
	"os/signal"
	"syscall"

	"github.com/phuslu/log"
)

func main() {
	_ = godotenv.Load()

	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     0,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

	log.Info().Msg("Starting LitePay...")

	r := router.Init()
	cryptocurrency.Init()

	go r.Listen(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")), fiber.ListenConfig{
		ListenerNetwork:       "tcp",
		DisableStartupMessage: true,
	})

	sChan := make(chan os.Signal, 1)
	defer close(sChan)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("LitePay started successfully!")

	<-sChan

	log.Info().Msg("Shutting down LitePay...")

	if err := r.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown LitePay.")
	}

	db.Close()
}
