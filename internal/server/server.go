package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

// Server implements a http.Hander and provides all the functionality
// required for creating new short links.
type Server struct {
	router  chi.Router
	options Options
}

// Options is a collection of settings exposed to configure a Server
// at generation time.
type Options struct {
	Logger      zerolog.Logger
	DB          *sql.DB
	ValidTokens []string
}

// WithLogger is a configurator that injects the given Logger into the
// final configuration.
func WithLogger(logger zerolog.Logger) Configurator {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithDatabase is a configurator that injects the given database
// instance into the final configuration.
func WithDatabase(db *sql.DB) Configurator {
	return func(o *Options) {
		o.DB = db
	}
}

type Configurator func(o *Options)

func (srv *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		srv.sendError(ctx, w, fmt.Errorf("no id found in the URL"), http.StatusBadRequest)
		return
	}
	tx, err := srv.options.DB.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		srv.sendError(ctx, w, fmt.Errorf("failed to open transaction: %w", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	var u string
	if err := tx.QueryRowContext(ctx, "SELECT url FROM links WHERE id = ?", id).Scan(&u); err != nil {
		srv.sendError(ctx, w, fmt.Errorf("not found: %w", err), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (srv *Server) handleCreateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := r.FormValue("url")
	if u == "" {
		srv.sendError(ctx, w, fmt.Errorf("url parameter missing"), http.StatusBadRequest)
		return
	}
	parsedu, err := url.Parse(u)
	if err != nil {
		srv.sendError(ctx, w, err, http.StatusBadRequest)
		return
	}
	tx, err := srv.options.DB.BeginTx(ctx, nil)
	if err != nil {
		srv.sendError(ctx, w, err, http.StatusInternalServerError)
		return
	}
	var foundID string
	err = tx.QueryRowContext(ctx, "SELECT id FROM links WHERE url = ?", parsedu.String()).Scan(&foundID)
	if err != nil && sql.ErrNoRows != err {
		tx.Rollback()
		srv.sendError(ctx, w, err, http.StatusInternalServerError)
		return
	}
	if foundID == "" {
		knownIDs, err := srv.getKnownIDs(ctx, tx)
		if err != nil {
			tx.Rollback()
			srv.sendError(ctx, w, err, http.StatusInternalServerError)
			return
		}
		foundID, err = srv.generateID(ctx, 5, knownIDs)
		if _, err := tx.ExecContext(ctx, "INSERT INTO links (id, url) VALUES(?, ?)", foundID, parsedu.String()); err != nil {
			tx.Rollback()
			srv.sendError(ctx, w, err, http.StatusInternalServerError)
			return
		}
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			srv.sendError(ctx, w, err, http.StatusInternalServerError)
			return
		}
	} else {
		tx.Rollback()
	}
	w.Write([]byte(foundID))
}

func (srv *Server) getKnownIDs(ctx context.Context, tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(ctx, "SELECT id FROM links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]string, 0, 10)
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

// generateID generates a new unique ID that is not present in the
// knownIDs. It tries to generate one with the preferred length. If
// that doesn't succeed after 10 attempts it will increase the length
// by 1.
func (srv *Server) generateID(ctx context.Context, preferredLength int, knownIDs []string) (string, error) {
	existing := make(map[string]struct{})
	for _, id := range knownIDs {
		existing[id] = struct{}{}
	}
	chars := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	length := preferredLength
	for {
		for i := 0; i < 10; i++ {
			candidate := &bytes.Buffer{}
			for j := 0; j < length; j++ {
				num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
				if err != nil {
					return "", err
				}
				candidate.WriteRune(chars[num.Int64()])
			}
			if _, ok := existing[candidate.String()]; !ok {
				return candidate.String(), nil
			}
		}
		length++
	}
}

// New creates a new Server instance after applying all provided
// configurators.
func New(configurators ...Configurator) *Server {
	options := Options{
		Logger: zerolog.Nop().With().Logger(),
	}
	for _, configurator := range configurators {
		configurator(&options)
	}
	srv := &Server{}
	router := chi.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(options.Logger.WithContext(ctx)))
		})
	})
	router.With(srv.tokenRequiredMiddleware).Post("/", srv.handleCreateURL)
	router.Get("/{id}", srv.handleResolve)
	srv.router = router
	srv.options = options
	return srv
}

func (srv *Server) tokenRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		if !srv.isValidToken(ctx, token) {
			srv.sendError(ctx, w, fmt.Errorf("no valid token provided"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (srv *Server) isValidToken(ctx context.Context, token string) bool {
	for _, v := range srv.options.ValidTokens {
		if fmt.Sprintf("SimpleToken %s", v) == token {
			return true
		}
	}
	return false

}

// ServeHTTP provided to implement the http.Handler interface.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *Server) sendError(ctx context.Context, w http.ResponseWriter, err error, status int) {
	logger := zerolog.Ctx(ctx)
	if status >= 500 {
		logger.Error().Err(err).Msg(err.Error())
	}
	http.Error(w, "", status)
}
