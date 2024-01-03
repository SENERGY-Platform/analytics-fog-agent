package agent

import (
	"encoding/json"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

func (agent *Agent) StopOperator(command operatorEntities.StopOperatorAgentControlCommand) {
	err := agent.ContainerManager.StopOperator(command.DeploymentReference)
	response := operatorEntities.OperatorAgentResponse{}
	response.Agent = conf.GetConf()
	response.OperatorId = command.OperatorID

	if err != nil {
		response.Success = false
		response.Error = err.Error()
		response.OperatorState = "not stopped"

	} else {
		response.Success = true
		response.Error = ""
		response.OperatorState = "stopped"
	}

	out, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	agent.PublishMessage(operatorEntities.StopOperatorResponseFogTopic, string(out), 2)
}

func (agent *Agent) StartOperator(command operatorEntities.StartOperatorControlCommand) {
	containerId, err := agent.ContainerManager.StartOperator(command)
	var responseMessage []byte
	if err != nil {
		// TODO add logging
		response := operatorEntities.StopOperatorAgentResponse{}
		response.Success = false
		response.Error = err.Error()
		response.OperatorState = "not started"
		response.Agent = conf.GetConf()
		response.OperatorId = command.Config.OperatorId
		responseMessage, err = json.Marshal(response)
		if err != nil {
			panic(err)
		}
	} else {
		response := operatorEntities.StartOperatorAgentResponse{}
		response.Success = true
		response.Error = ""
		response.OperatorState = "started"
		response.ContainerId = containerId
		response.OperatorId = command.Config.OperatorId
		response.Agent = conf.GetConf()
		responseMessage, err = json.Marshal(response)
		if err != nil {
			panic(err)
		}
	}

	agent.PublishMessage(operatorEntities.StartOperatorResponseFogTopic, string(responseMessage), 2)
}
