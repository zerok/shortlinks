package server_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/golang-migrate/migrate/v4"
	driver "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zerok/shortlinks/internal/server"
)

func setupDatabase(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	d, err := driver.WithInstance(db, &driver.Config{})
	require.NoError(t, err)
	mig, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", "../../migrations"), ":memory:", d)
	require.NoError(t, err)
	if err := mig.Up(); err != nil {
		if err != migrate.ErrNoChange {
			require.NoError(t, err)
		}
	}
	return db
}

type testbed struct {
	srv        *server.Server
	db         *sql.DB
	logger     zerolog.Logger
	validToken string
}

func (tb *testbed) teardown() {
	if tb.db != nil {
		tb.db.Close()
	}
}

func setupTestbed(t *testing.T) *testbed {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	db := setupDatabase(t)
	srv := server.New(func(c *server.Options) {
		c.ValidTokens = []string{"token"}
		c.DB = db
		c.Logger = logger
	})
	return &testbed{
		srv:        srv,
		logger:     logger,
		db:         db,
		validToken: "token",
	}
}

func TestCreateLinkMissingURL(t *testing.T) {
	tb := setupTestbed(t)
	defer tb.teardown()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("url="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "SimpleToken token")
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateLinkBrokenDB(t *testing.T) {
	tb := setupTestbed(t)
	// Close the DB right away
	tb.db.Close()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("url=http://test.com"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "SimpleToken token")
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateLink(t *testing.T) {
	tb := setupTestbed(t)
	defer tb.teardown()

	// Sending a POST request must require a valid token:
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusForbidden, w.Code)

	// Once a valid token is provided, a new link is stored and
	// its code returned:
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("url=http://test.com"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "SimpleToken token")
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	code := w.Body.String()
	require.NotEmpty(t, code)

	// That code should be usable:
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/"+code, nil)
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusTemporaryRedirect, w.Code)
	require.Equal(t, "http://test.com", w.Header().Get("Location"))

	// Creating a link with an already existing URL should return the existing code:
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("url=http://test.com"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "SimpleToken token")
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, code, w.Body.String())

	// Creating another link with an unknown URL should return a new code:
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("url=http://test2.com"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "SimpleToken token")
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.NotEmpty(t, w.Body.String())
	require.NotEqual(t, code, w.Body.String())
}

func TestResolveDBClosed(t *testing.T) {
	tb := setupTestbed(t)
	tb.teardown()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/something", nil)
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestResolveMissing(t *testing.T) {
	tb := setupTestbed(t)
	defer tb.teardown()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/something", nil)
	tb.srv.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}
