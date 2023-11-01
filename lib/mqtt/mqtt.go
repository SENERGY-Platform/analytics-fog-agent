package mqtt

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/topic"
)

func NewMQTTClient(brokerConfig config.BrokerConfig) *mqtt.MQTTClient {
	topics := mqtt.TopicConfig{
		topic.TopicPrefix + conf.GetConf().Id: byte(2),
		constants.MasterTopic:                 byte(2),
	}

	return &mqtt.MQTTClient{
		Broker:      brokerConfig,
		TopicConfig: topics,
	}
}
