package base

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func GetGormConfig() *gorm.Config {
	return &gorm.Config{
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
	}
}

type logWriter struct{}

func (w *logWriter) Printf(format string, data ...interface{}) {
	log.Infof(format, data...)
}
