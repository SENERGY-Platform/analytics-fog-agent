/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package relay

import (
	"fmt"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/agent"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	agentLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type RelayController struct {
	Agent *agent.Agent
}

func NewRelayController(agent *agent.Agent) *RelayController {
	return &RelayController{
		Agent: agent,
	}
}

func (relay *RelayController) ProcessMessage(message MQTT.Message) {
	switch message.Topic() {
	case constants.MasterTopic:
		relay.processMasterMessages(message.Payload())
	case agentLib.AgentsTopic + "/" + relay.Agent.Conf.Id:
		relay.processComandForAgent(message.Payload())
	}
}

func (relay *RelayController) OnMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	go relay.ProcessMessage(message)
}
