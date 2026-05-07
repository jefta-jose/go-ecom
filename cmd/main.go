package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	// uncomment and run go mod tidy incase you want to run migrations locally
	// "github.com/golang-migrate/migrate/v4"
    // _ "github.com/golang-migrate/migrate/v4/database/postgres"
    // _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// safely load the env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing...")
	}

	//-----------------------------------------------------------------
	// Run database migrations manually before starting the server
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

	// for running migrations locally
	// m, err := migrate.New(
	// 	"file:/home/devndegwa/Development/ecom-go-api-project/internal/adapters/postgresql/migrations",
	// 	"postgres://postgres:postgres123@localhost:5432/go_ecom?sslmode=disable",
	// )

	// if err != nil {
    //     log.Fatal(err)
    // }

    // if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    //     log.Fatal(err)
    // }

    // log.Println("Migrations applied successfully")

	api := application{
		config: cfg,
		db:     conn,
	}
	
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
