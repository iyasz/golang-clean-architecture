package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}

func NewDatabase(config *Database, log *logrus.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.Host,
		config.User,
		config.Password,
		config.Name,
		config.Port,
		config.Timezone,
	)

	log.Debugf(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold: 5 * time.Second,
			Colorful: false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries: true,
			LogLevel: logger.Info,
		}),	
	})

	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}

	connection.SetMaxIdleConns(config.Pool.Idle)
	connection.SetMaxOpenConns(config.Pool.Max)
	connection.SetConnMaxLifetime(time.Second * time.Duration(config.Pool.Lifetime))

	return db
}