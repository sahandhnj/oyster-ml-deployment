package node

import (
	"github.com/google/uuid"
)

type Node struct {
	ID          string `json:"id"`
	ImageName   string `json:"image_name"`
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

func NewNode(name string, pearlId string, containerId string, state State) (error, *Node) {
	n := Node{uuid.New().String(), name, pearlId, containerId, state}

	// json, err := json.Marshal(n)
	// jsonString := string(json[:])

	// writeToFiles(jsonString, ".oyster/"+n.ID+".json")

	// if err != nil {
	// 	return err, nil
	// }

	return nil, &n
}
