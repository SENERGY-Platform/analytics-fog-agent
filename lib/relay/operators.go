package relay 

import (
	"encoding/json"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

func (relay *RelayController) processStopOperatorCommand(message []byte) {
	logging.Logger.Debug("Received stop operator message")

	command := operatorEntities.StopOperatorAgentControlCommand{}
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
