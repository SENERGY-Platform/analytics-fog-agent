package container_manager

import (
	"errors"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/entities"
)

type Manager interface {
	StartOperator(operatorJob entities.OperatorJob) (containerId string, err error)
	StopOperator(operatorJob entities.OperatorJob) (err error)
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
