package container_manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"

	mqtt "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	aux_client "github.com/SENERGY-Platform/mgw-module-manager/aux-client"
	mgw_model "github.com/SENERGY-Platform/mgw-module-manager/lib/model"
)

type MGWManager struct {
	Broker        mqtt.BrokerConfig
	AuxDeploymentClient *aux_client.Client
	DeploymentID string
	ForcePullImg bool
}

func NewMGWManager(brokerHost, brokerPort, moduleManagerUrl, deploymentID string) *MGWManager {
	baseClient := &http.Client{}
	return &MGWManager{
		Broker: mqtt.BrokerConfig{
			Host: brokerHost,
			Port: brokerPort,
		},
		AuxDeploymentClient: aux_client.New(baseClient, moduleManagerUrl),
		DeploymentID: deploymentID,
		ForcePullImg: true,
	}
}
func (manager *MGWManager) CreateAndStartOperator(ctx context.Context, startRequest operatorEntities.StartOperatorControlCommand) (containerId string, err error) {
	createAuxRequest, err := manager.CreateAuxDeploymentRequest(startRequest)
	logging.Logger.Debug(fmt.Sprintf("Try to create aux deployment %s", createAuxRequest.Name))
	jobID, err := manager.AuxDeploymentClient.CreateAuxDeployment(ctx, manager.DeploymentID, createAuxRequest, manager.ForcePullImg)
	if err != nil {
		return "", err
	}
	logging.Logger.Debug(fmt.Sprintf("Wait for create aux deployment job %s", jobID))
	createJobResponse, err := manager.WaitForJob(ctx, jobID)
	if err != nil {
		return "", err
	}
	auxDeploymentID := createJobResponse.(string)
	logging.Logger.Debug(fmt.Sprintf("Try to start aux deployment %s", createAuxRequest.Name))
	jobID, err = manager.AuxDeploymentClient.StartAuxDeployment(ctx, manager.DeploymentID, auxDeploymentID)
	if err != nil {
		return "", err
	}
	logging.Logger.Debug(fmt.Sprintf("Wait for start aux deployment job %s", jobID))
	_, err = manager.WaitForJob(ctx, jobID)
	if err != nil {
		return "", err
	}
	return auxDeploymentID, nil
}

func (manager *MGWManager) RemoveOperator(ctx context.Context, operatorID string) (err error) {
	logging.Logger.Debug(fmt.Sprintf("Try to stop aux deployment %s", operatorID))
	jobID, err := manager.AuxDeploymentClient.StopAuxDeployment(ctx, manager.DeploymentID, operatorID)
	if err != nil {
		return err
	}
	logging.Logger.Debug(fmt.Sprintf("Wait for stop job %s", jobID))
	_, err = manager.WaitForJob(ctx, jobID)
	if err != nil {
		return err
	}
	logging.Logger.Debug(fmt.Sprintf("Try to remove aux deployment %s", operatorID))
	jobID, err = manager.AuxDeploymentClient.DeleteAuxDeployment(ctx, manager.DeploymentID, operatorID, false)
	if err != nil {
		return err
	}
	logging.Logger.Debug(fmt.Sprintf("Wait for remove job %s", jobID))
	_, err = manager.WaitForJob(ctx, jobID)
	return err
}

func (manager *MGWManager) GetOperatorState(ctx context.Context, operatorID string) (state OperatorState, err error) {
	auxDeployment, err := manager.AuxDeploymentClient.GetAuxDeployment(ctx, manager.DeploymentID, operatorID, false, true)
	if err != nil {
		return
	}
	state = OperatorState{
		State: auxDeployment.Container.Info.State,
	}
	return
}

func (manager *MGWManager) UpdateOperator(ctx context.Context, operatorID string, updateRequest operatorEntities.StartOperatorControlCommand) (err error) {
	updateAuxRequest, err := manager.CreateAuxDeploymentRequest(updateRequest)
	if err != nil {
		return err
	}
	jobID, err := manager.AuxDeploymentClient.UpdateAuxDeployment(ctx, manager.DeploymentID, operatorID, updateAuxRequest, false, manager.ForcePullImg)
	if err != nil {
		return err
	}
	_, err = manager.WaitForJob(ctx, jobID)
	return
}

func (manager *MGWManager) GetOperatorStates(ctx context.Context) (states map[string]OperatorState, err error) {
	filter := mgw_model.AuxDepFilter{}
	auxDeployments, err := manager.AuxDeploymentClient.GetAuxDeployments(ctx, manager.DeploymentID, filter, false, true)
	if err != nil {
		return
	}
	states = map[string]OperatorState{}
	for auxDepId, auxDeployment := range auxDeployments {
		states[auxDepId] = OperatorState{
			State: auxDeployment.Container.Info.State,
		}
	}
	return
}

func (manager *MGWManager) RestartOperator(ctx context.Context, operatorID string) (err error) {
	jobID, err := manager.AuxDeploymentClient.RestartAuxDeployment(ctx, manager.DeploymentID, operatorID)
	if err != nil {
		return err
	}
	_, err = manager.WaitForJob(ctx, jobID)
	return
}

type Logger struct {}

func (l Logger) Error(arg ...any) {
	fmt.Println(arg[0].(string))
}

func (manager *MGWManager) WaitForJob(ctx context.Context, jobID string) (result any, err error) {
	delay := 5 * time.Second
	timeout := 30 * time.Second
	logger := Logger{}
	jobResponse, err := aux_client.AwaitJob(ctx, manager.AuxDeploymentClient, manager.DeploymentID, jobID, delay, timeout, logger)
	if err != nil {
		return nil, err
	}
	if jobResponse.Error != nil {
		return jobResponse.Result, errors.New(fmt.Sprintf("Error: %s", jobResponse.Error.Message))
	}
	return  jobResponse.Result, nil
}

func (manager *MGWManager) CreateAuxDeploymentRequest(request operatorEntities.StartOperatorControlCommand) (auxDepRequest mgw_model.AuxDepReq, err error) {
	operatorConfig, inputTopics, config, err := StartOperatorConfigsToString(request)
	if err != nil {
		return
	}
	configs := map[string]string{
		"INPUT": inputTopics,
		"CONFIG": config,
		"OPERATOR_CONFIG": operatorConfig,
		"BROKER_HOST": manager.Broker.Host,
		"BROKER_PORT": manager.Broker.Port,
	}
	auxDepRequest = mgw_model.AuxDepReq{
		Image: request.ImageId,
		Configs: configs,
		Ref: "operator",
	}
	return
}