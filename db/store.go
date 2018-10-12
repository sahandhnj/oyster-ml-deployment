package db

import (
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sahandhnj/apiclient/store/model"
)

const (
	databaseFileName = "oysterbox.db"
	DBStorePath      = "~/.oysterdb"
	TIMEOUT_SECONDS  = 1
)

type DBStore struct {
	path         string
	db           *bolt.DB
	ModelService model.Service
}

func NewDBStore(DBStorePath string) (*DBStore, error) {
	databasePath := path.Join(DBStorePath, databaseFileName)

	DBStore := &DBStore{
		path: databasePath,
	}

	return DBStore, nil
}

func (d *DBStore) Open() error {
	db, err := bolt.Open(d.path, 0600, &bolt.Options{Timeout: TIMEOUT_SECONDS * time.Second})
	if err != nil {
		return err
	}
	DBStore.db = db

	return DBStore.initServices()
}

func (d *DBStore) Close() error {
	if d.db != nil {
		return d.db.Close()
	}

	return nil
}

func (d *DBStore) initServices() error {
	modelDBService, err := model.NewService(d.db)
	if err != nil {
		return err
	}

	f.ModelService = modelDBService

	return nil
}
