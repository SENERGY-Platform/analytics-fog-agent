package container_manager

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
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
func (mangager *MGWManager) StartOperator(operatorJob operatorEntities.StartOperatorMessage) (containerId string, err error) {
	return "id", nil
}

func (mangager *MGWManager) StopOperator(operatorId string) (err error) {
	return nil
}
