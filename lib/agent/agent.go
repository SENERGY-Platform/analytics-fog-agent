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
	"fmt"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/container_manager"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/entities"
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
	conf, _ := json.Marshal(entities.AgentMessage{Type: "register", Conf: conf.GetConf()})
	agent.PublishMessage(constants.AgentsTopic, string(conf), 2)
}

func (agent *Agent) processMessage(message MQTT.Message) {
	command := entities.ControlCommand{}
	err := json.Unmarshal(message.Payload(), &command)
	if err != nil {
		fmt.Println("error:", err)
	}
	switch command.Command {
	case "startOperator":
		containerId, err := agent.ContainerManager.StartOperator(command.Data)
		if err != nil {
			command.Data.Response = "Error"
			command.Data.ResponseMessage = err.Error()
		} else {
			command.Data.Response = "OK"
			command.Data.ResponseMessage = "All good"
			command.Data.ContainerId = containerId
		}
		command.Data.Agent = conf.GetConf()
		out, err := json.Marshal(command.Data)
		if err != nil {
			panic(err)
		}

		agent.PublishMessage(constants.OperatorsTopic, string(out), 2)
	case "stopOperator":
		agent.ContainerManager.StopOperator(command.Data)
	case "ping":
		agent.sendPong()
	}

}

func (agent *Agent) sendPong() {
	out, err := json.Marshal(entities.AgentMessage{Type: "pong", Conf: conf.GetConf()})
	if err != nil {
		panic(err)
	}
	agent.PublishMessage(constants.AgentsTopic, string(out), 1)
}

func (agent *Agent) onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	go agent.processMessage(message)
}

func (agent *Agent) PublishMessage(topic string, message string, qos int) {
	agent.Client.Publish(topic, message, qos)
}
