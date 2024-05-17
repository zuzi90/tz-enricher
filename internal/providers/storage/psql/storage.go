package psql

import (
	"context"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.up.sql
var embedMigrations embed.FS

type Storage struct {
	db    *sqlx.DB
	dbDSN string
	log   *logrus.Entry
}

func NewStorage(log *logrus.Logger, dbDSN string) (*Storage, error) {
	stor := Storage{
		dbDSN: dbDSN,
		log:   log.WithField("module", "storage"),
	}

	clientSqlx, err := sqlx.Open("pgx", stor.dbDSN)
	if err != nil {
		return nil, fmt.Errorf("err client sqlx: %v", err)
	}

	stor.db = clientSqlx

	return &stor, nil
}

func (s *Storage) getMigrations() error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("discovering migrations by caller: %w", err)
	}

	goose.SetBaseFS(embedMigrations)
	err = goose.Up(s.db.DB, "migrations")
	if err != nil {
		return fmt.Errorf("error when trying to access the migration files: %w", err)
	}

	return nil

}
func (s *Storage) UpdateSchema() error {
	if err := s.getMigrations(); err != nil {
		return err
	}

	s.log.Info("Migration applied")

	return nil
}

func (s *Storage) CloseDB() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *Storage) TruncateTables(tables ...string) error {

	ctx := context.Background()

	for _, table := range tables {
		_, err := s.db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return err
		}
	}

	return nil
}
