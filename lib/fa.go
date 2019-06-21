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

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func RegisterAgent() {
	conf, _ := json.Marshal(AgentMessage{Type: "register", Conf: GetConf()})
	publishMessage("agents", string(conf))
}

func processMessage(message MQTT.Message) {
	command := ControlCommand{}
	err := json.Unmarshal(message.Payload(), &command)
	if err != nil {
		fmt.Println("error:", err)
	}
	switch command.Command {
	case "startOperator":
		containerId := startOperator(command.Data)
		command.Data.ContainerId = containerId
		command.Data.Agent = GetConf()
		out, err := json.Marshal(command.Data)
		if err != nil {
			panic(err)
		}
		publishMessage("operators", string(out))
	case "stopOperator":
		stopOperator(command.Data)
	case "ping":
		sendPong()
	}

}

func startOperator(operator OperatorJob) (containerId string) {
	PullImage(operator.ImageId)
	containerId = RunContainer(operator.ImageId)
	return
}

func stopOperator(operatorJob OperatorJob) {
	StopContainer(operatorJob.ContainerId)
	RemoveContainer(operatorJob.ContainerId)
}

func sendPong() {
	out, err := json.Marshal(AgentMessage{Type: "pong", Conf: GetConf()})
	if err != nil {
		panic(err)
	}
	publishMessage("agents", string(out))
}
