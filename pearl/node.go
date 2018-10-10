package pearl

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Node struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
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

func NewNode(pearlId string, name string, containerId string, state State) (error, *Node) {
	n := Node{uuid.New().String(), name, pearlId, containerId, state}

	json, err := json.Marshal(n)
	jsonString := string(json[:])

	writeToFiles(jsonString, ".oyster/"+n.ID+".json")

	if err != nil {
		return err, nil
	}

	return nil, &n
}
