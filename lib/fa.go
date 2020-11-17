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
package lib

import (
	"encoding/json"
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func RegisterAgent() {
	conf, _ := json.Marshal(AgentMessage{Type: "register", Conf: GetConf()})
	publishMessage(AgentsTopic, string(conf))
}

func processMessage(message MQTT.Message) {
	command := ControlCommand{}
	err := json.Unmarshal(message.Payload(), &command)
	if err != nil {
		fmt.Println("error:", err)
	}
	switch command.Command {
	case "startOperator":
		containerId, err := startOperator(command.Data)
		if err != nil {
			command.Data.Response = "Error"
			command.Data.ResponseMessage = err.Error()
		} else {
			command.Data.Response = "OK"
			command.Data.ResponseMessage = "All good"
			command.Data.ContainerId = containerId
		}
		command.Data.Agent = GetConf()
		out, err := json.Marshal(command.Data)
		if err != nil {
			panic(err)
		}
		publishMessage(OperatorsTopic, string(out))
	case "stopOperator":
		stopOperator(command.Data)
	case "ping":
		sendPong()
	}

}

func startOperator(operator OperatorJob) (containerId string, err error) {
	operatorConfig, err := json.Marshal(operator.OperatorConfig)
	if err != nil {
		panic(err)
	}
	inputTopics, err := json.Marshal(operator.InputTopics)
	if err != nil {
		panic(err)
	}
	config, err := json.Marshal(operator.Config)
	if err != nil {
		return
	}
	env := []string{
		"INPUT=" + string(inputTopics),
		"CONFIG=" + string(config),
		"OPERATOR_CONFIG=" + string(operatorConfig),
		"BROKER_HOST=" + GetEnv("CONTAINER_BROKER_HOST", GetEnv("BROKER_HOST", "localhost")),
		"BROKER_PORT=" + GetEnv("BROKER_PORT", "1883"),
	}
	pull, err := strconv.ParseBool(GetEnv("CONTAINER_PULL_IMAGE", "true"))
	if err != nil {
		fmt.Println("Invalid config for CONTAINER_PULL_IMAGE")
		pull = true
	}
	containerId, err = RunContainer(operator.ImageId, env, pull, operator.Config.PipelineId, operator.Config.OperatorId)
	return
}

func stopOperator(operatorJob OperatorJob) {
	RemoveContainer(operatorJob.ContainerId)
}

func sendPong() {
	out, err := json.Marshal(AgentMessage{Type: "pong", Conf: GetConf()})
	if err != nil {
		panic(err)
	}
	publishMessage(AgentsTopic, string(out))
}
