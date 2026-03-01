package postgres

import (
	"Assignment2/pkg/modules"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewDialect(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var (
		db  *sqlx.DB
		err error
	)

	for attempt := 1; attempt <= 10; attempt++ {
		db, err = sqlx.ConnectContext(ctx, "postgres", dsn)
		if err == nil {
			break
		}

		log.Printf("database is not ready yet (attempt %d/10): %v", attempt, err)
		if attempt < 10 {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		panic(fmt.Errorf("connect database: %w", err))
	}

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	log.Println("Database connection established")

	autoMigrate(cfg)
	log.Println("Database migrations completed")

	return &Dialect{DB: db}
}

func autoMigrate(cfg *modules.PostgreConfig) {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
