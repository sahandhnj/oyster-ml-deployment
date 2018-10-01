package pearl

import (
	"fmt"
	"time"
)

type Pearl struct {
	Name        string
	CreatedAt   time.Time
	Description string
	ModelPath   string
}

func NewPearl(name string, desc string, modelPath string) *Pearl {
	return &Pearl{name, time.Now(), desc, modelPath}
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
