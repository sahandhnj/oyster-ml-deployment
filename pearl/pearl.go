package pearl

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

type Pearl struct {
	ID          bson.ObjectId `bson:"_id" json:"id" yaml:"id"`
	Name        string        `bson:"name" json:"name" yaml:"name"`
	Description string        `bson:"description" json:"description" yaml:"description"`
	ModelPath   string        `bson:"model_path" json:"model_path" yaml:"model_path"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at" yaml:"created_at"`
}

func NewPearl(name string, desc string, modelPath string) (*Pearl, error) {
	p := Pearl{bson.NewObjectId(), name, desc, modelPath, time.Now()}

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

func (p *Pearl) deploy() {

}

func (p *Pearl) setup() {

}
