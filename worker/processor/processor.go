package processor

import (
	"encoding/json"
	"fmt"
)

var (
	processorMap = make(map[string]TaskProcessor)
)

type TaskRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type TaskProcessor interface {
	GetType() string
	Process(TaskRequest) error
}

func RegisterProcessor(p TaskProcessor) error {
	if _, ok := processorMap[p.GetType()]; ok {
		return fmt.Errorf("duplicate processor for type: %s", p.GetType())
	} else {
		processorMap[p.GetType()] = p
	}

	return nil
}

func Process(r TaskRequest) error {
	if p, ok := processorMap[r.Type]; ok {
		return p.Process(r)
	} else {
		return fmt.Errorf("invalid processor for type: %s", r.Type)
	}
}
