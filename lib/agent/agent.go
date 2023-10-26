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

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/container_manager"
	agentEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"
	controlEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Agent struct {
	ContainerManager container_manager.Manager
	Client           *MQTTClient
}

func NewAgent(containerManager container_manager.Manager, mqttClient *MQTTClient) *Agent {
	return &Agent{
		ContainerManager: containerManager,
		Client:           mqttClient,
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

func (agent *Agent) processMessage(message MQTT.Message) {
	command := controlEntities.ControlMessage{}
	err := json.Unmarshal(message.Payload(), &command)
	if err != nil {
		logging.Logger.Error("error:", err)
	}

	switch command.Command {
	case "startOperator":
		logging.Logger.Debug("Received start operator message")
		command := operatorEntities.StartOperatorControlCommand{}
		err := json.Unmarshal(message.Payload(), &command)
		if err != nil {
			logging.Logger.Error("error:", err)
		}

		containerId, err := agent.ContainerManager.StartOperator(command.Operator)
		var message []byte
		if err != nil {
			response := operatorEntities.OperatorAgentResponse{}
			response.Response = "Error"
			response.ResponseMessage = err.Error()
			response.Agent = conf.GetConf()
			response.OperatorId = command.Operator.Config.OperatorId
			message, err = json.Marshal(response)
			if err != nil {
				panic(err)
			}
		} else {
			response := operatorEntities.OperatorAgentSuccessResponse{}
			response.Response = "OK"
			response.ResponseMessage = "All good"
			response.ContainerId = containerId
			response.OperatorId = command.Operator.Config.OperatorId
			response.Agent = conf.GetConf()
			message, err = json.Marshal(response)
			if err != nil {
				panic(err)
			}
		}

		agent.PublishMessage(constants.OperatorsTopic, string(message), 2)
	case "stopOperator":
		logging.Logger.Debug("Received stop operator message")

		command := operatorEntities.StopOperatorControlCommand{}
		err := json.Unmarshal(message.Payload(), &command)
		if err != nil {
			logging.Logger.Error("error:", err)
		}

		err = agent.ContainerManager.StopOperator(command.OperatorId)
		response := operatorEntities.OperatorAgentResponse{}
		response.Agent = conf.GetConf()
		response.OperatorId = command.OperatorId

		if err != nil {
			response.Response = "Error"
			response.ResponseMessage = err.Error()

		} else {
			response.Response = "OK"
			response.ResponseMessage = "All good"
		}

		out, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}

		agent.PublishMessage(constants.OperatorsTopic, string(out), 2)
	case "ping":
		logging.Logger.Debug("Received ping message")
		agent.sendPong()
	}

}

func (agent *Agent) sendPong() {
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

func (agent *Agent) onMessageReceived(client MQTT.Client, message MQTT.Message) {
	logging.Logger.Debugf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	go agent.processMessage(message)
}

func (agent *Agent) PublishMessage(topic string, message string, qos int) {
	agent.Client.Publish(topic, message, qos)
}
