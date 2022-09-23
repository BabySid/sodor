package metastore

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func initGorm(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(&logWriter{}, logger.Config{
			SlowThreshold:             time.Second * 3,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Info,
		}),
		DryRun:      false,
		PrepareStmt: true,
		QueryFields: true,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

type logWriter struct{}

func (w *logWriter) Printf(format string, data ...interface{}) {
	log.Infof(format, data)
}
