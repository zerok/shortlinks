package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	driver "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/zerok/shortlinks/internal/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	var addr string
	var dbPath string
	var tokens []string
	pflag.StringVar(&addr, "addr", "localhost:8000", "Address to listen on")
	pflag.StringVar(&dbPath, "db", "./shortlinks.sqlite", "Path to a database file")
	pflag.StringSliceVar(&tokens, "token", []string{}, "Valid tokens for creating new links")
	pflag.Parse()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to open database at %s.", dbPath)
	}

	d, err := driver.WithInstance(db, &driver.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to initialize migration driver.")
	}

	mig, err := migrate.NewWithDatabaseInstance("file://./migrations", dbPath, d)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to initialize migrations.")
	}
	if err := mig.Up(); err != nil {
		if err != migrate.ErrNoChange {
			logger.Fatal().Err(err).Msgf("Failed to apply migrations.")
		}
	}

	handler := server.New(server.WithLogger(logger), server.WithDatabase(db), func(o *server.Options) {
		o.ValidTokens = tokens
	})

	s := http.Server{}
	s.Handler = handler
	s.Addr = addr
	logger.Info().Msgf("Listening on %s", addr)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start HTTP server.")
	}
}
