package asset

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/wutipong/albums2/worker/processor"
)

type Payload struct {
	AssetID string `json:"asset_id"`
}

type Processor struct{}

func (p *Processor) GetType() string {
	return "asset"
}

func (p *Processor) Process(req processor.TaskRequest) error {
	payload := Payload{}
	err := json.Unmarshal(req.Payload, &payload)
	if err != nil {
		return fmt.Errorf("unable to unmarshal the payload: %w", err)
	}

	slog.Info("processing asset", "id", payload.AssetID)

	return nil
}
