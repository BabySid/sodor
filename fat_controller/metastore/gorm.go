package metastore

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sodor/base"
	"time"
)

func initGorm(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), base.GetGormConfig())

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
