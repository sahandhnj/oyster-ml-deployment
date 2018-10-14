package version

import (
	"github.com/boltdb/bolt"
	"github.com/sahandhnj/apiclient/types"
	"github.com/sahandhnj/apiclient/util"
)

const (
	BucketName = "version"
)

type Service struct {
	db *bolt.DB
}

func NewService(db *bolt.DB) (*Service, error) {
	err := util.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

func (s *Service) Version(ID int) (*types.Version, error) {
	var version types.Version
	identifier := util.Itob(int(ID))

	err := util.GetObject(s.db, BucketName, identifier, &version)
	if err != nil {
		return nil, err
	}

	return &version, nil
}

func (s *Service) VersionsByModelId(modelId int) ([]types.Version, error) {
	var versions = make([]types.Version, 0)

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var ver types.Version
			err := util.UnmarshalJsonObject(v, &ver)
			if err != nil {
				return err
			}

			if ver.ModelID == modelId {
				versions = append(versions, ver)
			}
		}

		return nil
	})

	return versions, err
}

func (s *Service) VersionByVersionNumber(versionNumber int) (*types.Version, error) {
	var version *types.Version

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var t types.Version
			err := util.UnmarshalJsonObject(v, &t)
			if err != nil {
				return err
			}

			if t.VersionNumber == versionNumber {
				version = &t
				break
			}
		}

		if version == nil {
			return util.GetError(util.ErrNotFound)
		}

		return nil
	})

	return version, err
}

func (s *Service) VersionByName(name string) (*types.Version, error) {
	var version *types.Version

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var t types.Version
			err := util.UnmarshalJsonObject(v, &t)
			if err != nil {
				return err
			}

			if t.Name == name {
				version = &t
				break
			}
		}

		if version == nil {
			return util.GetError(util.ErrNotFound)
		}

		return nil
	})

	return version, err
}

func (s *Service) Versions() ([]types.Version, error) {
	var versions = make([]types.Version, 0)

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var version types.Version
			err := util.UnmarshalJsonObject(v, &version)
			if err != nil {
				return err
			}
			versions = append(versions, version)
		}

		return nil
	})

	return versions, err
}

func (s *Service) GetNextIdentifier() int {
	return util.GetNextIdentifier(s.db, BucketName)
}

func (s *Service) CreateVersion(version *types.Version) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		err := bucket.SetSequence(uint64(version.ID))
		if err != nil {
			return err
		}

		data, err := util.MarshalJsonObject(version)
		if err != nil {
			return err
		}

		return bucket.Put(util.Itob(int(version.ID)), data)
	})
}

func (s *Service) UpdateVersion(ID int, version *types.Version) error {
	identifier := util.Itob(int(ID))
	return util.UpdateObject(s.db, BucketName, identifier, version)
}

func (s *Service) DeleteVersion(ID int) error {
	identifier := util.Itob(int(ID))
	return util.DeleteObject(s.db, BucketName, identifier)
}
