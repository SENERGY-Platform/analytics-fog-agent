package agent

import (
	"encoding/json"
	"time"

	"context"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
)

func (agent *Agent) StopOperator(command operatorEntities.StopOperatorAgentControlCommand) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10 * time.Second)
	defer cancel()
	err := agent.ContainerManager.RemoveOperator(ctx, command.DeploymentReference)
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
	ctx, cancel := context.WithTimeout(context.TODO(), 10 * time.Second)
	defer cancel()
	logging.Logger.Debug("Try to start operator: " + command.ImageId)
	containerId, err := agent.ContainerManager.CreateAndStartOperator(ctx, command)
	var responseMessage []byte
	if err != nil {
		response := operatorEntities.StartOperatorAgentResponse{}
		response.Success = false
		response.Error = err.Error()
		response.OperatorState = "not started"
		response.Agent = conf.GetConf()
		response.OperatorId = command.Config.OperatorId
		responseMessage, err = json.Marshal(response)
		if err != nil {
			logging.Logger.Errorf("Could not unmarshal response %s", err.Error())
		}
		logging.Logger.Errorf("Could not start operator %s", response.Error)
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
			logging.Logger.Errorf("Could not unmarshal response %s", err.Error())
		} else {
			logging.Logger.Infof("Operator started successfully: %s", responseMessage)
		}
	}

	agent.PublishMessage(operatorEntities.StartOperatorResponseFogTopic, string(responseMessage), 2)
}
