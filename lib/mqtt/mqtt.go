package mqtt

import (
	"fmt"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
)

func NewMQTTClient(brokerConfig mqtt.BrokerConfig, logger *slog.Logger) *mqtt.MQTTClient {
	agentID := conf.GetConf().Id
	topics := mqtt.TopicConfig{
		agent.AgentsTopic + "/" + agentID: byte(2),
		agent.GetStartOperatorAgentTopic(agentID): byte(2),
		agent.GetStopOperatorAgentTopic(agentID): byte(2),
		constants.MasterTopic:                       byte(2),
	}

	return &mqtt.MQTTClient{
		Broker:      brokerConfig,
		TopicConfig: topics,
		Logger:      logger,
		OnConnectHandler: OnConnect,
		SubscribeInitial: true,
	}
}

func OnConnect(client MQTT.Client) {
	fmt.Println("MQTT client connected!")
}