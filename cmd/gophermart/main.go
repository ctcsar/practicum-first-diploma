package main

import (
	"context"
	"database/sql"
	"flag"

	chi "github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"

	f "github.com/ctcsar/practicum-first-diploma/internal/flags"
	"github.com/ctcsar/practicum-first-diploma/internal/handlers"
	"github.com/ctcsar/practicum-first-diploma/internal/logger"
)

func main() {
	ctx := context.Background()
	flags := f.NewServerFlags()
	flags.SetServerFlags()
	flag.Parse()
	handler := chi.NewRouter()
	db, err := sql.Open("pgx", flags.GetDatabasePath())
	if err != nil {
		logger.Log.Fatal("cannot connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Log.Fatal("cannot ping database", zap.Error(err))
	}
	err = goose.Up(db, "../../internal/database/migrations")
	if err != nil {
		logger.Log.Fatal("cannot migrate database", zap.Error(err))
	}
	err = handlers.Run(ctx, flags.GetServerURL(), handler, db)
	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}
}
