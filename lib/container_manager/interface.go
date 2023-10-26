package container_manager

import (
	"errors"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

type Manager interface {
	StartOperator(operatorJob operatorEntities.StartOperatorMessage) (containerId string, err error)
	StopOperator(operatorID string) (err error)
}

func NewManager(config *config.Config) (Manager, error) {
	switch config.ContainerManager {
	case constants.DockerManager:
		return NewDockerManager(config.Broker.Host, config.Broker.Port, config.ContainerPullImage, config.ContainerNetwork), nil
	case constants.MGWManager:
		return NewMGWManager(config.Broker.Host, config.Broker.Port, config.ModuleManager.Host, config.ModuleManager.Port), nil
	}

	return nil, errors.New("Container Manager not found")
}
