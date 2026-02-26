package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// load dotenv
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MigrateDIR := os.Getenv("MigrateDIR")
	GoLangMigrateURL := os.Getenv("GoLangMigrateURL")
	
	if GoLangMigrateURL == "" || MigrateDIR == "" {
		log.Fatal("Environment variables GoLangMigrateURL and MigrateDIR must be set")
	}
	//-----------------------------------------------------------------

	// run migrations
	m, err := migrate.New(
		"file://"+MigrateDIR,
		GoLangMigrateURL,
	)

    if err != nil {
        log.Fatal(err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migrations applied successfully")
	//-----------------------------------------------------------------

	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: os.Getenv("DBSTRING"),
		},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)

	if err != nil {
		panic(err)
	}

	defer conn.Close(ctx)

	logger.Info("connected to database", "", "")

	api := application{
		config: cfg,
		db:     conn,
	}
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
