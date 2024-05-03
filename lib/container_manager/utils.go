package container_manager

import (
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"encoding/json"
)

func StartOperatorConfigsToString(startRequest operatorEntities.StartOperatorControlCommand) (string, string, string, error) {
	operatorConfig, err := json.Marshal(startRequest.OperatorConfig)
	if err != nil {
		return "","","", err
	}
	inputTopics, err := json.Marshal(startRequest.InputTopics)
	if err != nil {
		return "","","", err
	}
	config, err := json.Marshal(startRequest.Config)
	if err != nil {
		return "","","", err
	}
	return string(operatorConfig), string(inputTopics), string(config), nil
}