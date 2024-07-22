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
	"context"
	"encoding/json"
	"time"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib"
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
	StorageHandler lib.StorageHandler
}

func NewAgent(containerManager container_manager.Manager, mqttClient *mqtt.MQTTClient, conf agentEntities.Configuration, controlOperatorTimeout time.Duration, storageHandler lib.StorageHandler) *Agent {
	return &Agent{
		ContainerManager: containerManager,
		Client:           mqttClient,
		Conf:             conf,
		ControlOperatorTimeout: controlOperatorTimeout,
		StorageHandler: storageHandler,
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
		CurrentOperatorStates: []agentEntities.OperatorState{},
	})
	agent.PublishMessage(constants.AgentsTopic, string(conf), 2)
}

func (agent *Agent) SendPong() {
	ctx := context.Background()
	allOperateStates, err := agent.StorageHandler.GetOperatorStates(ctx)
	if err != nil {
		logging.Logger.Error("Cant load all operator states during pong", "error", err.Error())
		return
	}
	out, err := json.Marshal(agentEntities.AgentInfoMessage{
		ControlMessage: controlEntities.ControlMessage{
			Command: "pong",
		},
		Conf: conf.GetConf(),
		CurrentOperatorStates: allOperateStates,
	})
	if err != nil {
		logging.Logger.Error("Cant marshal pong message", "error", err.Error())
	}
	agent.PublishMessage(constants.AgentsTopic, string(out), 1)
}

func (agent *Agent) PublishMessage(topic string, message string, qos int) {
	logging.Logger.Debug("Publish message", "message", message, "topic", topic, "qos", qos)
	agent.Client.Publish(topic, message, qos)
}
