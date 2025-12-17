package database

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/harryosmar/protobuf-go/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates and returns a new database connection with connection pooling and retry logic
func NewDatabase(cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DatabaseConnectTimeout)*time.Second)
	defer cancel()

	return NewDatabaseWithContext(ctx, cfg, zapLogger)
}

// NewDatabaseWithContext creates a database connection with context support and retry logic
func NewDatabaseWithContext(ctx context.Context, cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
	// Configure GORM logger to use Zap
	gormLogger := logger.New(
		&GormZapWriter{logger: zapLogger},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection with retry logic
	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= cfg.DatabaseMaxRetries; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("database connection cancelled: %w", ctx.Err())
		default:
		}

		zapLogger.Info("Attempting database connection",
			zap.Int("attempt", attempt),
			zap.Int("max_retries", cfg.DatabaseMaxRetries),
		)

		db, err = gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{
			Logger: gormLogger,
		})

		if err == nil {
			// Connection successful, test it
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				pingErr := sqlDB.Ping()
				if pingErr == nil {
					break // Success
				}
				err = pingErr
			} else {
				err = dbErr
			}
		}

		// Connection failed
		if attempt < cfg.DatabaseMaxRetries {
			// Calculate exponential backoff delay
			delay := time.Duration(cfg.DatabaseRetryDelay) * time.Second * time.Duration(math.Pow(2, float64(attempt-1)))
			zapLogger.Warn("Database connection failed, retrying",
				zap.Error(err),
				zap.Int("attempt", attempt),
				zap.Duration("retry_in", delay),
			)

			// Wait with context awareness
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("database connection cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
				// Continue to next attempt
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", cfg.DatabaseMaxRetries, err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool for high-traffic
	sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DatabaseMaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DatabaseMaxLife) * time.Second)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	zapLogger.Info("Database connected successfully",
		zap.String("max_idle", fmt.Sprintf("%d", cfg.DatabaseMaxIdle)),
		zap.String("max_open", fmt.Sprintf("%d", cfg.DatabaseMaxOpen)),
		zap.String("max_lifetime", fmt.Sprintf("%ds", cfg.DatabaseMaxLife)),
	)

	return db, nil
}

// GormZapWriter implements GORM's logger interface using Zap
type GormZapWriter struct {
	logger *zap.Logger
}

func (g *GormZapWriter) Printf(format string, args ...interface{}) {
	g.logger.Info(fmt.Sprintf(format, args...))
}

// CloseDatabase closes the database connection
func CloseDatabase(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
