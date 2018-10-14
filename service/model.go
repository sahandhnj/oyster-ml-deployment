package service

import (
	"errors"

	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/filemanager"
	"github.com/sahandhnj/apiclient/types"
	"github.com/sahandhnj/apiclient/util"
)

type ModelService struct {
	Model     *types.Model
	Nodes     []*types.Node `json:"nodes"`
	file      *filemanager.FileStoreManager
	DBHandler *db.DBStore
}

func NewModelService(model *types.Model, dbHandler *db.DBStore) (*ModelService, error) {
	file, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	modelService := ModelService{
		Nodes:     nil,
		file:      file,
		DBHandler: dbHandler,
	}

	configExists, err := file.ConfigFileExists()
	if err != nil {
		return nil, err
	}

	if configExists {
		model, err = ReadModel()
		if err != nil {
			return nil, err
		}
	} else {
		if model == nil {
			return nil, errors.New("No oyster project found. use init command to start one")
		}

		model.ID = dbHandler.ModelService.GetNextIdentifier()
		dbHandler.ModelService.CreateModel(model)
	}

	modelService.Model = model

	err = modelService.file.WriteToConfigFile(model)

	if err != nil {
		return nil, err
	}

	return &modelService, nil
}

func ReadModel() (*types.Model, error) {
	file, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	data, err := file.ReadConfigFile()
	if err != nil {
		return nil, err
	}

	var model types.Model
	err = util.UnmarshalYamlObject(data, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
