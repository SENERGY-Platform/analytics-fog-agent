package container_manager

import (
	"errors"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"context"
)

type Manager interface {
	CreateAndStartOperator(ctx context.Context, operatorJob operatorEntities.StartOperatorControlCommand) (containerId string, err error)
	RestartOperator(ctx context.Context, operatorID string) (err error)
	RemoveOperator(ctx context.Context, operatorID string) (err error)
	GetOperatorState(ctx context.Context, operatorID string) (state OperatorState, err error)
	UpdateOperator(ctx context.Context, operatorID string, updateRequest operatorEntities.StartOperatorControlCommand) (err error)
	GetOperatorStates(ctx context.Context, ) (states map[string]OperatorState, err error)
}

func NewManager(config *config.Config) (Manager, error) {
	switch config.ContainerManager {
	case constants.DockerManager:
		// Agent might not run as container and could have localhost as hostname
		return NewDockerManager(config.ContainerBrokerHost, config.ContainerBrokerPort, config.ContainerPullImage, config.ContainerNetwork), nil
	case constants.MGWManager:
		return NewMGWManager(config.Broker.Host, config.Broker.Port, config.ModuleManagerURL, config.DeploymentID), nil
	}

	return nil, errors.New("Container Manager not found")
}
