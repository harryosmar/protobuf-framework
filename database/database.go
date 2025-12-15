package database

import (
	"fmt"
	"time"

	"github.com/harryosmar/protobuf-go/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates and returns a new database connection with connection pooling
func NewDatabase(cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
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

	// Open database connection
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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
