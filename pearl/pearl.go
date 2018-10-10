package pearl

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Pearl struct {
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	ModelPath   string    `json:"model_path" yaml:"model_path"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
	Nodes       []*Node   `json:"nodes"`
}

func NewPearl(name string, desc string, modelPath string) (*Pearl, error) {
	p := Pearl{uuid.New().String(), name, desc, modelPath, time.Now(), nil}

	yamlBytes, err := yaml.Marshal(p)
	yamlString := string(yamlBytes[:])
	writeToFiles(yamlString, "oyster.yml")

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func ReadPearl() (*Pearl, error) {
	data, err := ioutil.ReadFile("oyster.yml")
	if err != nil {
		log.Fatal(err)
	}

	var pearl Pearl
	if err := pearl.Parse(data); err != nil {
		log.Fatal(err)
	}

	return &pearl, nil
}

func (p *Pearl) UpdateConfig() error {
	yamlBytes, err := yaml.Marshal(p)
	yamlString := string(yamlBytes[:])
	writeToFiles(yamlString, "oyster.yml")

	return err
}

func (p *Pearl) Parse(data []byte) error {
	return yaml.Unmarshal(data, p)
}

func (p *Pearl) Config() {
	makeDockerFile()
}

func (p *Pearl) PrintInfo() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Model Path: %s\n", p.ModelPath)
}
