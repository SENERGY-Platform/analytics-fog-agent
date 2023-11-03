package container_manager

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	mqtt "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

type MGWManager struct {
	ModuleManager config.ModuleManagerConfig
	Broker        mqtt.BrokerConfig
}

func NewMGWManager(brokerHost string, brokerPort string, moduleManagerHost string, moduleManagerPort string) *MGWManager {
	return &MGWManager{
		Broker: mqtt.BrokerConfig{
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

	//
	return "id", nil
}

func (mangager *MGWManager) StopOperator(operatorId string) (err error) {
	return nil
}
