package model

import (
	"fmt"
)

type Model struct {
	ID          int    `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	ModelPath   string `json:"model_path" yaml:"model_path"`
}

func NewModel(ID int, name string, desc string, modelPath string) *Model {
	model := Model{
		ID:          ID,
		Name:        name,
		Description: desc,
		ModelPath:   modelPath,
	}

	return &model
}

func (m *Model) PrintInfo() {
	fmt.Printf("Name: %s\n", m.Name)
	fmt.Printf("Description: %s\n", m.Description)
	fmt.Printf("Model Path: %s\n", m.ModelPath)
}
