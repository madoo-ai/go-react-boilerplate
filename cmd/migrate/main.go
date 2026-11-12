package main

import (
	"flag"
	"fmt"
	"os"

	pgadapter "go-react-boilerplate/internal/adapter/postgres"
	"go-react-boilerplate/internal/xmodule"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "migrate:", err)
		os.Exit(1)
	}
}

func run() error {
	dsn := flag.String("dsn", "", "Postgres DSN (env DATABASE_URL)")
	flag.Parse()

	cfg := xmodule.LoadConfig()
	if *dsn != "" {
		cfg.DSN = *dsn
	}
	if cfg.DSN == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if err := pgadapter.Migrate(cfg.DSN); err != nil {
		return err
	}
	fmt.Println("Migrations applied.")
	return nil
}
