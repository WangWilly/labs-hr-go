package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

////////////////////////////////////////////////////////////////////////////////

type DbConfig struct {
	Host     string `env:"DB_HOST,default=localhost"`
	Port     string `env:"DB_PORT,default=3306"`
	User     string `env:"DB_USER,default=labs-hr-go"`
	Password string `env:"DB_PASSWORD,default=labs-hr-go"`
	Database string `env:"DB_DATABASE,default=labs-hr-go"`
	Charset  string `env:"DB_CHARSET,default=utf8mb4"`
	Driver   string `env:"DB_DRIVER,default=mysql"`
	Timezone string `env:"DB_TIMEZONE,default=UTC"`

	SlowThreshold time.Duration `env:"DB_SLOW_THRESHOLD,default=200ms"`
	IsDev         bool          `env:"DB_IS_DEV,default=true"`
}

const (
	// DBContextKey is the key used to store the database connection in the context
	DBContextKey contextKey = "gormdb"
)

////////////////////////////////////////////////////////////////////////////////

func GetDB(cfg DbConfig) (*gorm.DB, error) {
	////////////////////////////////////////////////////////////////////////////
	logLevel := glogger.Silent
	if cfg.IsDev {
		logLevel = glogger.Info
	}

	logger := glogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), glogger.Config{
		SlowThreshold:             cfg.SlowThreshold,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	////////////////////////////////////////////////////////////////////////////
	// Initialize the database connection

	var db *gorm.DB
	var err error

	switch cfg.Driver {
	case "mysql":
		dsn := cfg.User + ":" + cfg.Password +
			"@tcp(" + cfg.Host + ":" + cfg.Port + ")/" +
			cfg.Database + "?charset=" + cfg.Charset + "&parseTime=true"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   logger,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			return nil, err
		}
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	////////////////////////////////////////////////////////////////////////////
	// Set the connection pool settings

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, err
	// }
	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// // Set the connection timeout
	// sqlDB.SetConnMaxLifetime(time.Hour)
	// sqlDB.SetConnMaxIdleTime(time.Hour)

	////////////////////////////////////////////////////////////////////////////

	return db, nil
}

func GetCtxDb(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(DBContextKey).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to get database from context")
	}
	return db, nil
}
