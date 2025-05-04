package manager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver

	commonConfig "github.com/zydee3/stockdb/internal/common/config"
	"github.com/zydee3/stockdb/internal/common/logger"
)

const (
	databaseSavePath = commonConfig.BaseSavePath + "/manager_jobs.db"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	// Check if the database file already exists, if not, create it
	if _, err := os.Stat(databaseSavePath); os.IsNotExist(err) {
		if createError := os.MkdirAll(commonConfig.BaseSavePath, 0750); createError != nil {
			return nil, fmt.Errorf("failed to create directory: %w", createError)
		}
		if createError := os.WriteFile(databaseSavePath, []byte{}, 0600); createError != nil {
			return nil, fmt.Errorf("failed to create database file: %w", createError)
		}
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf(
		"file:%s?_journal=WAL&_sync=NORMAL&_timeout=5000&_fk=true",
		databaseSavePath,
	))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if initError := createJobSchema(db); initError != nil {
		return nil, fmt.Errorf("%w", initError)
	}

	logger.Infof("Manager database running on %s", databaseSavePath)

	return &Storage{db: db}, nil
}

func createJobSchema(db *sql.DB) error {
	const (
		timeoutValue    = 5
		timeoutDuration = time.Duration(timeoutValue) * time.Second
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackError := tx.Rollback(); rollbackError != nil && !errors.Is(rollbackError, sql.ErrTxDone) {
			logger.Errorf("%v", rollbackError)
		}
	}()

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id TEXT NOT NULL UNIQUE,
			job_type TEXT NOT NULL CHECK(job_type IN ('RECURRING', 'INTERVAL')),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			schedule_type TEXT NOT NULL CHECK(schedule_type IN ('RECURRING', 'INTERVAL')),
			schedule_start_at DATETIME NOT NULL,
			schedule_end_at DATETIME,
			schedule_frequency TEXT,
			spec_json TEXT NOT NULL CHECK(json_valid(spec_json)),
			attempts INTEGER NOT NULL DEFAULT 0 CHECK(attempts >= 0),
			max_retries INTEGER NOT NULL DEFAULT 3 CHECK(max_retries >= 0)
		)
	`)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = tx.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_realtime 
		ON jobs (created_at DESC) 
		WHERE date(created_at) = date('now')
	`)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = tx.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_historical 
		ON jobs (created_at ASC) 
		WHERE date(created_at) < date('now')
	`)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return tx.Commit()
}
