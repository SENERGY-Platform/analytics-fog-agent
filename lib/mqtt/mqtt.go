package mqtt

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	log_level "github.com/y-du/go-log-level"
)

func NewMQTTClient(brokerConfig mqtt.BrokerConfig, logger *log_level.Logger) *mqtt.MQTTClient {
	topics := mqtt.TopicConfig{
		agent.AgentsTopic + "/" + conf.GetConf().Id: byte(2),
		constants.MasterTopic:                       byte(2),
	}

	return &mqtt.MQTTClient{
		Broker:      brokerConfig,
		TopicConfig: topics,
		Logger:      logger,
	}
}
