package metastore

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
)

func parseDSN(uri string) (string, error) {
	scheme := "mysql://"
	if !strings.HasPrefix(uri, scheme) {
		return "", fmt.Errorf("dsn must start with '%s'", scheme)
	}
	return uri[len(scheme):], nil
}

type metaStore struct {
	db *gorm.DB
}

var (
	once      sync.Once
	singleton *metaStore
	URI       string
)

func GetInstance() *metaStore {
	once.Do(func() {
		singleton = &metaStore{}
		err := singleton.initOnce(URI)
		if err != nil {
			log.Fatalf("metastore init failed. err=%s", err)
		}
	})
	return singleton
}

func (ms *metaStore) initOnce(uri string) error {
	if ms.db != nil {
		return nil
	}

	dsn, err := parseDSN(uri)
	if err != nil {
		return err
	}

	db, err := initGorm(dsn)
	if err != nil {
		return err
	}

	ms.db = db
	return nil
}

func (ms *metaStore) AutoMigrate() error {
	db := ms.db.Set("gorm:table_options", "ENGINE=InnoDB")
	for _, tbl := range totalTable {
		if db.Migrator().HasTable(tbl) {
			err := db.Migrator().DropTable(tbl)
			if err != nil {
				log.Warn(err)
				return err
			}
		}
	}
	for _, tbl := range totalTable {
		err := db.Migrator().CreateTable(tbl)
		if err != nil {
			log.Warn(err)
			return err
		}
	}
	return nil
}
