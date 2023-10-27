package relay

import (
	"encoding/json"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
	controlEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	masterEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/master"
)

func (relay *RelayController) processMasterMessages(message []byte) {
	command := controlEntities.ControlMessage{}
	err := json.Unmarshal(message, &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	switch command.Command {
	case "register":
		relay.processMasterRegisterMessage(message)
	}
}

func (relay *RelayController) processMasterRegisterMessage(message []byte) {
	logging.Logger.Debug("Received master register message")

	command := masterEntities.MasterInfoMessage{}
	err := json.Unmarshal(message, &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	relay.Agent.Register()
}
