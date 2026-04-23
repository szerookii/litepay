package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/cron"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/cryptocurrency/bitcoin"
	"github.com/szerookii/litepay/backend/cryptocurrency/litecoin"
	"github.com/szerookii/litepay/backend/cryptocurrency/solana"
	"github.com/szerookii/litepay/backend/db"
	"github.com/szerookii/litepay/backend/router"
	"github.com/szerookii/litepay/backend/secrets"
	"github.com/szerookii/litepay/backend/utils/env"
)

func main() {
	_ = godotenv.Load()

	if err := env.Check(
		env.WithDefault("HOST", "0.0.0.0"),
		env.WithDefault("PORT", "8080"),
		env.Required("DATABASE_URL"),
		env.Required("JWT_SECRET", env.MinLen(32)),
		env.Optional("SECRET_PROVIDER", env.OneOf("env", "vault", "bitwarden", "aws", "gcp")),
		env.Optional("ALLOW_REGISTER", env.OneOf("true", "false")),
		env.Optional("BTC_RPC_URL"),
		env.Optional("LTC_RPC_URL"),
		env.Optional("SOL_RPC_URL"),
		env.Optional("ALLOWED_ORIGINS"),
	); err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err := secrets.Load(context.Background()); err != nil {
		log.Fatal().Msg(err.Error())
	}

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

	cryptocurrency.Register(&bitcoin.Bitcoin{})
	cryptocurrency.Register(&litecoin.Litecoin{})
	cryptocurrency.Register(&solana.Solana{})

	r := router.Init()

	ctxVerifier, cancelVerifier := context.WithCancel(context.Background())
	go cron.StartTransactionVerifier(ctxVerifier, 30*time.Second)

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info().Msgf("Started web server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Web server failed")
		}
	}()

	sChan := make(chan os.Signal, 1)
	defer close(sChan)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("LitePay started successfully!")
	<-sChan

	log.Info().Msg("Shutting down LitePay...")
	cancelVerifier()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown LitePay.")
	}
	db.Close()
}
