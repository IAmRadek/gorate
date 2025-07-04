package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IAmRadek/go-kit/envconfig"
	"github.com/IAmRadek/gorate/internal/exchanges"
	"github.com/IAmRadek/gorate/internal/rates"
	"github.com/gin-gonic/gin"
)

type Config struct {
	GinMode                  string        `env:"GIN_MODE" default:"debug"`
	Addr                     string        `env:"ADDR" default:":8080"`
	ReadTimeout              time.Duration `env:"READ_TIMEOUT" default:"10s"`
	ReadHeaderTimeout        time.Duration `env:"READ_HEADER_TIMEOUT" default:"10s"`
	WriteTimeout             time.Duration `env:"WRITE_TIMEOUT" default:"10s"`
	IdleTimeout              time.Duration `env:"IDLE_TIMEOUT" default:"10s"`
	MaxHeaderBytes           int           `env:"MAX_HEADER_BYTES" default:"1024"`
	GracefulShutdownDuration time.Duration `env:"GRACEFUL_SHUTDOWN_DURATION" default:"5s"`

	OpenExchangeRatesProviderAppID string `env:"OPEN_EXCHANGE_RATES_PROVIDER_APP_ID" required:"true"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log := slog.Default()

	var cfg Config
	err := envconfig.Read(&cfg, os.LookupEnv)
	if err != nil {
		fatal("reading config: %v", err)
	}

	httpClient := http.DefaultClient

	ratesProvider := rates.NewOpenExchangeRatesProvider(httpClient, cfg.OpenExchangeRatesProviderAppID)

	fixedCryptoRates := rates.NewFixedCryptoRatesProvider()
	exchange := exchanges.NewExchange(fixedCryptoRates)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	registerRoutes(router, ratesProvider, exchange)

	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	log.Info("Starting Server", "addr", cfg.Addr)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("Server Failed", "err", err)
			}
		}
	}()

	<-ctx.Done()

	log.Info("Shutting down server...")
	teardownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownDuration)
	defer cancel()

	if err := httpSrv.Shutdown(teardownCtx); err != nil {
		log.Error("Server Failed to Shutdown", "err", err)
	}

	log.Info("Server Stopped")
}

func registerRoutes(
	router *gin.Engine,
	provider *rates.OpenExchangeRatesProvider,
	exchange *exchanges.Exchange,
) {
	router.GET("/rates", HandleRates(provider))
	router.GET("/exchange", HandleExchange(exchange))
}

func fatal(msg string, a ...any) {
	_, _ = fmt.Fprintf(os.Stderr, msg, a...)
	os.Exit(-1)
}
