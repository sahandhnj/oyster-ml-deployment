package model

import (
	"fmt"

	"github.com/sahandhnj/apiclient/types"
	"github.com/sahandhnj/apiclient/types/node"
	"github.com/sahandhnj/apiclient/util"

	"github.com/google/uuid"
	"github.com/sahandhnj/apiclient/filemanager"
)

type Model struct {
	Config *types.ModelConfig
	Nodes  []*node.Node `json:"nodes"`
	file   *filemanager.FileStoreManager
}

func NewModel(name string, desc string, modelPath string) (*Model, error) {
	file, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	configExists, err := file.ConfigFileExists()
	if err != nil {
		return nil, err
	}
	if configExists {
		return ReadModel()
	}

	model := Model{
		Config: &types.ModelConfig{
			ID:          uuid.New().String(),
			Name:        name,
			Description: desc,
			ModelPath:   modelPath,
		},
		Nodes: nil,
		file:  file,
	}
	err = model.file.WriteToConfigFile(model.Config)

	if err != nil {
		return nil, err
	}

	return &model, nil
}

func ReadModel() (*Model, error) {
	file, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	data, err := file.ReadConfigFile()
	if err != nil {
		return nil, err
	}

	var config types.ModelConfig
	err = util.UnmarshalYamlObject(data, &config)
	if err != nil {
		return nil, err
	}

	model := Model{
		Config: &config,
		Nodes:  nil,
		file:   file,
	}

	return &model, nil
}

// func (p *Model) Config() {
// 	makeDockerFile()
// }

func (m *Model) PrintInfo() {
	fmt.Printf("Name: %s\n", m.Config.Name)
	fmt.Printf("Description: %s\n", m.Config.Description)
	fmt.Printf("Model Path: %s\n", m.Config.ModelPath)
}
