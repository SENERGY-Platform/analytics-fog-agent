package container_manager

import (
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"encoding/json"
)

type Config struct {
	OutputTopic string `json:"outputTopic"`
	PipelineID string `json:"pipelineId"`
	OperatorID string `json:"operatorId"`
	BaseOperatorID string `json:"baseOperatorId"`
} 

func StartOperatorConfigsToString(startRequest operatorEntities.StartOperatorControlCommand) (string, string, string, error) {
	operatorConfig, err := json.Marshal(startRequest.OperatorConfig)
	if err != nil {
		return "","","", err
	}
	inputTopics, err := json.Marshal(startRequest.InputTopics)
	if err != nil {
		return "","","", err
	}
	config := Config{
		OutputTopic: startRequest.OutputTopic,
		OperatorID: startRequest.OperatorIDs.OperatorId,
		PipelineID: startRequest.OperatorIDs.PipelineId,
		BaseOperatorID: startRequest.OperatorIDs.BaseOperatorId,
	}
	configMarshaled, err := json.Marshal(config)
	if err != nil {
		return "","","", err
	}
	return string(operatorConfig), string(inputTopics), string(configMarshaled), nil
}