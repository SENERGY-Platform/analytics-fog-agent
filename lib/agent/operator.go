package agent

import (
	"encoding/json"
	"time"

	"context"

	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
)

func (agent *Agent) StopOperator(command operatorEntities.StopOperatorAgentControlCommand) {
	ctx, cancel := context.WithTimeout(context.TODO(), agent.ControlOperatorTimeout * time.Second)
	defer cancel()
	err := agent.StorageHandler.SaveOperatorState(ctx, command.PipelineId, command.OperatorId, "stopping", "", "", nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state", "error", err.Error())
		return
	}
	err = agent.ContainerManager.RemoveOperator(ctx, command.ContainerId)
	if err != nil {
		logging.Logger.Error("Could not remove operator", "operatorID", command.OperatorId, "error", err)
	}
	response := operatorEntities.OperatorAgentResponse{}
	response.AgentId = agent.Conf.Id
	response.OperatorId = command.OperatorId
	response.Time = time.Now()
	response.PipelineId = command.PipelineId

	if err != nil {
		response.Error = err.Error()
		response.DeploymentState = "not stopped"
		err = agent.StorageHandler.SaveOperatorState(ctx, command.PipelineId, command.OperatorId, response.DeploymentState, "", "", nil)
		if err != nil {
			logging.Logger.Error("Could not save starting state", "error", err.Error())
			return
		}
	} else {
		response.Error = ""
		response.DeploymentState = "stopped"
		err = agent.StorageHandler.DeleteOperator(ctx, command.PipelineId, command.OperatorId, nil)
		if err != nil {
			logging.Logger.Error("Could not delete stopped operator", "error", err.Error())
			return
		}
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
	logging.Logger.Debug("Try to start operator: " + command.ImageId, "pipelineID", command.PipelineId, "operatorID", command.OperatorId)
	err := agent.StorageHandler.SaveOperatorState(ctx, command.OperatorIDs.PipelineId, command.OperatorIDs.OperatorId, "starting", "", "", nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state", "error", err.Error())
		return
	}

	containerId, err := agent.ContainerManager.CreateAndStartOperator(ctx, command)
	var responseMessage []byte
	response := operatorEntities.StartOperatorAgentResponse{}
	response.Time = time.Now()
	response.AgentId = agent.Conf.Id
	response.OperatorId = command.OperatorId
	response.PipelineId = command.PipelineId
	if err != nil {
		response.Error = err.Error()
		response.DeploymentState = "not started"
		responseMessage, err = json.Marshal(response)
		if err != nil {
			logging.Logger.Error("Could not unmarshal response", "error", err.Error())
		}
		logging.Logger.Error("Could not start operator", "error", response.Error)
	} else {
		response.Error = ""
		response.DeploymentState = "started"
		response.ContainerId = containerId
		responseMessage, err = json.Marshal(response)
		if err != nil {
			logging.Logger.Error("Could not unmarshal response", "error", err.Error())
		} 
		logging.Logger.Info("Operator started successfully: " + string(responseMessage))
	}

	err = agent.StorageHandler.SaveOperatorState(ctx, command.OperatorIDs.PipelineId, command.OperatorIDs.OperatorId, response.DeploymentState, response.ContainerId, response.Error, nil)
	if err != nil {
		logging.Logger.Error("Could not save starting state", "error", err.Error())
		return
	}

	agent.PublishMessage(operatorEntities.StartOperatorResponseFogTopic, string(responseMessage), 2)
}
