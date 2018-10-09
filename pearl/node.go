package pearl

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Node struct {
	ID          string `json:"id"`
	PearlId     string `json:"pearl_id"`
	ContainerId string `json:"container_id"`
	State       State  `json:"state"`
}

type State string

const (
	AnyState = State("")
	Running  = State("running")
	Stopped  = State("stopped")
)

func NewNode(pearlId string, containerId string, state State) (*Node, error) {
	n := Node{uuid.New().String(), pearlId, containerId, state}

	json, err := json.Marshal(n)
	jsonString := string(json[:])

	writeToFiles(jsonString, ".oyster/"+n.ID+".json")

	if err != nil {
		return nil, err
	}

	return &n, nil
}
