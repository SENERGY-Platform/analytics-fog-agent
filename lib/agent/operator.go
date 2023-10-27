package agent

import (
	"encoding/json"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

func (agent *Agent) StopOperator(command operatorEntities.StopOperatorControlCommand) {
	err := agent.ContainerManager.StopOperator(command.OperatorId)
	response := operatorEntities.OperatorAgentResponse{}
	response.Agent = conf.GetConf()
	response.OperatorId = command.OperatorId

	if err != nil {
		response.Response = operatorEntities.OperatorDeployedError
		response.ResponseMessage = err.Error()

	} else {
		response.Response = operatorEntities.OperatorDeployedSuccessfully
		response.ResponseMessage = "All good"
	}

	out, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	agent.PublishMessage(constants.OperatorsTopic, string(out), 2)
}

func (agent *Agent) StartOperator(command operatorEntities.StartOperatorControlCommand) {
	containerId, err := agent.ContainerManager.StartOperator(command.Operator)
	var responseMessage []byte
	if err != nil {
		response := operatorEntities.OperatorAgentResponse{}
		response.Response = "Error"
		response.ResponseMessage = err.Error()
		response.Agent = conf.GetConf()
		response.OperatorId = command.Operator.Config.OperatorId
		responseMessage, err = json.Marshal(response)
		if err != nil {
			panic(err)
		}
	} else {
		response := operatorEntities.OperatorAgentSuccessResponse{}
		response.Response = "OK"
		response.ResponseMessage = "All good"
		response.ContainerId = containerId
		response.OperatorId = command.Operator.Config.OperatorId
		response.Agent = conf.GetConf()
		responseMessage, err = json.Marshal(response)
		if err != nil {
			panic(err)
		}
	}

	agent.PublishMessage(constants.OperatorsTopic, string(responseMessage), 2)
}
