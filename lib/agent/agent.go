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
package agent

import (
	"encoding/json"
	"time"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/container_manager"
	agentEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"
	controlEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
)

type Agent struct {
	ContainerManager container_manager.Manager
	Client           *mqtt.MQTTClient
	Conf             agentEntities.Configuration
	ControlOperatorTimeout time.Duration
}

func NewAgent(containerManager container_manager.Manager, mqttClient *mqtt.MQTTClient, conf agentEntities.Configuration, controlOperatorTimeout time.Duration) *Agent {
	return &Agent{
		ContainerManager: containerManager,
		Client:           mqttClient,
		Conf:             conf,
		ControlOperatorTimeout: controlOperatorTimeout,
	}
}

func (agent *Agent) Register() {
	agentConf := conf.GetConf()
	logging.Logger.Debug("Register agent ", agentConf.Id)
	conf, _ := json.Marshal(agentEntities.AgentInfoMessage{
		ControlMessage: controlEntities.ControlMessage{
			Command: "register",
		},
		Conf: agentConf,
	})
	agent.PublishMessage(constants.AgentsTopic, string(conf), 2)
}

func (agent *Agent) SendPong() {
	out, err := json.Marshal(agentEntities.AgentInfoMessage{
		ControlMessage: controlEntities.ControlMessage{
			Command: "pong",
		},
		Conf: conf.GetConf(),
	})
	if err != nil {
		panic(err)
	}
	agent.PublishMessage(constants.AgentsTopic, string(out), 1)
}

func (agent *Agent) PublishMessage(topic string, message string, qos int) {
	agent.Client.Publish(topic, message, qos)
}
