package db

import (
	"fmt"
	"time"

	"github.com/juanfont/headscale/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// HSDatabase is the main database struct for headscale.
type HSDatabase struct {
	DB *gorm.DB
}

// NewHeadscaleDatabase initializes and returns a new HSDatabase instance
// based on the provided configuration. Supports SQLite and PostgreSQL.
func NewHeadscaleDatabase(cfg config.DatabaseConfig) (*HSDatabase, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	switch cfg.Type {
	case "sqlite", "sqlite3":
		if cfg.Sqlite.Path == "" {
			return nil, fmt.Errorf("sqlite database path is not set")
		}
		db, err = gorm.Open(sqlite.Open(cfg.Sqlite.Path), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}

		// Enable WAL mode for better concurrent read performance
		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			return nil, fmt.Errorf("failed to get underlying sql.DB: %w", sqlErr)
		}
		if _, sqlErr = sqlDB.Exec("PRAGMA journal_mode=WAL;"); sqlErr != nil {
			return nil, fmt.Errorf("failed to set WAL mode: %w", sqlErr)
		}

	case "postgres", "postgresql":
		dsn := buildPostgresDSN(cfg.Postgres)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres database: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &HSDatabase{DB: db}, nil
}

// buildPostgresDSN constructs a PostgreSQL DSN string from config.
func buildPostgresDSN(cfg config.PostgresConfig) string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Pass,
		cfg.Name,
	)
	if cfg.Ssl == "disable" || cfg.Ssl == "false" {
		dsn += " sslmode=disable"
	} else {
		dsn += fmt.Sprintf(" sslmode=%s", cfg.Ssl)
	}
	return dsn
}

// Ping verifies the database connection is alive.
func (h *HSDatabase) Ping() error {
	sqlDB, err := h.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}

// Close closes the underlying database connection.
func (h *HSDatabase) Close() error {
	sqlDB, err := h.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
