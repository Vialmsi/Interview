package store

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/Vialmsi/Interview/internal/config"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

const (
	//PgUniqueEntryErrorCode postgres code for check unique violation
	PgUniqueEntryErrorCode = "23505"
)

type Store struct {
	db           *sql.DB
	ctxTimeout   time.Duration
	defaultLimit int
	defaultPage  int
}

func NewStore(cfg config.PSQLDatabase) (*Store, error) {
	connStr := cfg.Address
	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("couldn't open PSQL: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("couldn't ping PSQl: %s", err)
	}

	logrus.Info("Ping PSQL - OK!")

	err = migrateTables(db)
	if err != nil {
		return nil, fmt.Errorf("couldn't migrate tables: %s", err)
	}

	return &Store{
		db:           db,
		ctxTimeout:   time.Second * time.Duration(cfg.Timeout),
		defaultLimit: cfg.DefaultLimit,
		defaultPage:  cfg.DefaultPage,
	}, nil
}

//go:embed migrations/tables/*.sql
var embedTablesMigrations embed.FS

func migrateTables(db *sql.DB) error {

	goose.SetBaseFS(embedTablesMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, "migrations/tables", goose.WithAllowMissing())

	return err
}

func isUniqueViolation(err error) bool {
	pgErr, ok := err.(*pq.Error)
	return ok && pgErr.Code == PgUniqueEntryErrorCode
}
