package container_manager

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/entities"
)

type MGWManager struct {
	ModuleManager config.ModuleManagerConfig
	Broker        config.BrokerConfig
}

func NewMGWManager(brokerHost string, brokerPort string, moduleManagerHost string, moduleManagerPort string) *MGWManager {
	return &MGWManager{
		Broker: config.BrokerConfig{
			Host: brokerHost,
			Port: brokerPort,
		},
		ModuleManager: config.ModuleManagerConfig{
			Host: moduleManagerHost,
			Port: moduleManagerPort,
		},
	}
}
func (mangager *MGWManager) StartOperator(operatorJob entities.OperatorJob) (containerId string, err error) {
	return "", nil
}

func (mangager *MGWManager) StopOperator(operatorJob entities.OperatorJob) (err error) {
	return nil
}
