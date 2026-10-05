package database

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
)

type Initializer struct {
	pgx *pgx.Tx
}

func RunMigrations(dbUrl string, migrationsDir string) error {
	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("resolve migrations dir %q: %w", migrationsDir, err)
	}

	sourceUrl := "file://" + filepath.ToSlash(absDir)
	migrations, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		return fmt.Errorf("init migrations: %w", err)
	}
	if err := migrations.Up(); err != nil {
		defer func() {
			soureErr, dbErr := migrations.Close()
			_ = soureErr
			_ = dbErr
		}()
	}

	if err = migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
