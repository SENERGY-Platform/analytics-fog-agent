package relay

import (
	"encoding/json"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
	controlEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

func (relay *RelayController) processComandForAgent(message []byte) {
	command := controlEntities.ControlMessage{}
	err := json.Unmarshal(message, &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	switch command.Command {
	case "startOperator":
		relay.processStartOperatorCommand(message)
	case "stopOperator":
		relay.processStopOperatorCommand(message)
	case "ping":
		logging.Logger.Debug("Received ping message")
		relay.Agent.SendPong()
	}
}

func (relay *RelayController) processStopOperatorCommand(message []byte) {
	logging.Logger.Debug("Received stop operator message")

	command := operatorEntities.StopOperatorControlCommand{}
	err := json.Unmarshal(message, &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	relay.Agent.StopOperator(command)
}

func (relay *RelayController) processStartOperatorCommand(message []byte) {
	logging.Logger.Debug("Received start operator message")
	command := operatorEntities.StartOperatorControlCommand{}
	err := json.Unmarshal(message, &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	relay.Agent.StartOperator(command)
}
