package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	pgadapter "go-react-boilerplate/internal/adapter/postgres"
	"go-react-boilerplate/internal/adapter/security"
	"go-react-boilerplate/internal/xmodule"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "seed:", err)
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
		return fmt.Errorf("migrate: %w", err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		TranslateError: true,
		Logger:         gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	hash, err := security.NewHasher().Hash("password")
	if err != nil {
		return fmt.Errorf("hash demo password: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	result, err := pgadapter.Seed(ctx, db, hash)
	if err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	fmt.Printf("Seed complete: %d permissions, %d roles, %d new users.\n",
		result.Permissions, result.Roles, result.NewUsers)
	return nil
}
