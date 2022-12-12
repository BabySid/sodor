package routine

import (
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"sodor/thomas/config"
	"time"
)

type removeTaskInstance struct {
}

func (r removeTaskInstance) Run() {
	dirs := make([]string, 0)
	now := time.Now()
	root := config.GetInstance().DataPath
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			if now.Sub(info.ModTime()) > config.GetInstance().InstanceMaxAge {
				dirs = append(dirs, path)
			}
		}

		return nil
	})

	if err != nil {
		log.Warnf("filepath.WalkDir return err=%s", err)
	}

	for _, fp := range dirs {
		_ = os.RemoveAll(fp)
		log.Infof("removeTaskInstance %s", fp)
	}
}
