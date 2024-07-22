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
	ctx, cancel := context.WithTimeout(context.TODO(), agent.ControlOperatorTimeout * time.Second)
	defer cancel()
	err := agent.StorageHandler.SaveOperatorState(ctx, command.PipelineID, command.OperatorID, "stopping", "", "", nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state %s", err.Error())
		return
	}
	err = agent.ContainerManager.RemoveOperator(ctx, command.DeploymentReference)
	if err != nil {
		logging.Logger.Error("Could not remove operator %s: %s", command.OperatorID, err.Error())
	}
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

	err = agent.StorageHandler.SaveOperatorState(ctx, command.PipelineID, command.OperatorID, response.OperatorState, "", "", nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state", "error", err.Error())
		return
	}

	out, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	agent.PublishMessage(operatorEntities.StopOperatorResponseFogTopic, string(out), 2)
}

func (agent *Agent) StartOperator(command operatorEntities.StartOperatorControlCommand) {
	ctx, cancel := context.WithTimeout(context.TODO(), 60 * time.Second)
	defer cancel()
	logging.Logger.Debug("Try to start operator: " + command.ImageId)
	err := agent.StorageHandler.SaveOperatorState(ctx, command.Config.OperatorIDs.PipelineId, command.Config.OperatorIDs.OperatorId, "starting", "", "", nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state", "error", err.Error())
		return
	}

	containerId, err := agent.ContainerManager.CreateAndStartOperator(ctx, command)
	var responseMessage []byte
	response := operatorEntities.StartOperatorAgentResponse{}
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		response.OperatorState = "not started"
		response.Agent = conf.GetConf()
		response.OperatorId = command.Config.OperatorId
		responseMessage, err = json.Marshal(response)
		if err != nil {
			logging.Logger.Error("Could not unmarshal response", "error", err.Error())
		}
		logging.Logger.Error("Could not start operator", "error", response.Error)
	} else {
		response.Success = true
		response.Error = ""
		response.OperatorState = "started"
		response.ContainerId = containerId
		response.OperatorId = command.Config.OperatorId
		response.Agent = conf.GetConf()
		responseMessage, err = json.Marshal(response)
		if err != nil {
			logging.Logger.Error("Could not unmarshal response", "error", err.Error())
		} 
		logging.Logger.Info("Operator started successfully: " + string(responseMessage))
	}

	err = agent.StorageHandler.SaveOperatorState(ctx, command.Config.OperatorIDs.PipelineId, command.Config.OperatorIDs.OperatorId, response.OperatorState, response.ContainerId, response.Error, nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state", "error", err.Error())
		return
	}

	agent.PublishMessage(operatorEntities.StartOperatorResponseFogTopic, string(responseMessage), 2)
}
