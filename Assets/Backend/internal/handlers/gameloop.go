package handlers

import (
	"encoding/json"

	"github.com/MisterDodik/MultiplayerGame/internal/events"
	"github.com/MisterDodik/MultiplayerGame/internal/network"
)

type PositionUpdatePayload struct {
}

func UpdatePositionHandler(e events.Event, c *network.Client) error {
	var payload PositionUpdatePayload
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return err
	}

	return nil
}
