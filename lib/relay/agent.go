package relay

import (
	"encoding/json"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
	controlEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
)

func (relay *RelayController) processComandForAgent(message []byte) {
	command := controlEntities.ControlMessage{}
	err := json.Unmarshal(message, &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	switch command.Command {
	case "ping":
		logging.Logger.Debug("Received ping message")
		relay.Agent.SendPong()
	}
}

