package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-react-boilerplate/internal/xmodule"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "server:", err)
		os.Exit(1)
	}
}

func run() error {
	listen := flag.String("listen", "", "listen address (env LISTEN, default 127.0.0.1:8080)")
	dsn := flag.String("dsn", "", "Postgres DSN (env DATABASE_URL)")
	flag.Parse()

	cfg := xmodule.LoadConfig()
	if *listen != "" {
		cfg.Listen = *listen
	}
	if *dsn != "" {
		cfg.DSN = *dsn
	}
	cfg.ResolveCookieSecure()

	secret, err := xmodule.LoadOrCreateSecret(cfg.AuthSecret)
	if err != nil {
		return fmt.Errorf("auth secret: %w", err)
	}
	cfg.AuthSecret = secret

	log := xmodule.NewLogger(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	deps, err := xmodule.InitDependencies(ctx, cfg, log)
	if err != nil {
		return err
	}
	return xmodule.Serve(ctx, cfg, deps, log)
}
