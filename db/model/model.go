package model

import (
	"github.com/boltdb/bolt"
	types "github.com/sahandhnj/apiclient/types/model"
	"github.com/sahandhnj/apiclient/util"
)

const (
	BucketName = "project"
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

func (s *Service) Model(ID int) (*types.Model, error) {
	var model types.Model
	identifier := util.Itob(int(ID))

	err := util.GetObject(Service.db, BucketName, identifier, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *Service) ModelByName(name string) (*types.Model, error) {
	var model *types.Model

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var t types.Model
			err := util.UnmarshalJsonObject(v, &t)
			if err != nil {
				return err
			}

			if t.Name == name {
				model = &t
				break
			}
		}

		if model == nil {
			return util.GetError(ErrNotFound)
		}

		return nil
	})

	return model, err
}

func (s *Service) Models() ([]types.Model, error) {
	var models = make([]types.Model, 0)

	err := Service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var Model types.Model
			err := util.UnmarshalJsonObject(v, &Model)
			if err != nil {
				return err
			}
			models = append(models, Model)
		}

		return nil
	})

	return models, err
}

func (s *Service) GetNextIdentifier() int {
	return util.GetNextIdentifier(s.db, BucketName)
}

func (s *Service) CreateModel(model *types.Model) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		err := bucket.SetSequence(uint64(model.ID))
		if err != nil {
			return err
		}

		data, err := util.MarshalJsonObject(model)
		if err != nil {
			return err
		}

		return bucket.Put(util.Itob(int(model.ID)), data)
	})
}

func (s *Service) UpdateModel(ID types.ModelID, model *types.Model) error {
	identifier := util.Itob(int(ID))
	return util.UpdateObject(s.db, BucketName, identifier, model)
}

func (s *Service) DeleteModel(ID types.ModelID) error {
	identifier := util.Itob(int(ID))
	return util.DeleteObject(Service.db, BucketName, identifier)
}
