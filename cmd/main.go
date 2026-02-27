package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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

	api := application{
		config: cfg,
		db:     conn,
	}
	
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
