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

package main

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/agent"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/container_manager"
	srv_base "github.com/SENERGY-Platform/go-service-base/srv-base"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	config, err := config.NewConfig("")
	if err != nil {
		log.Print("Error loading config")
	}
	log.Println("config: %s", srv_base.ToJsonStr(config))

	conf.InitConf()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	mqttClient := agent.NewMQTTClient(config.Broker)

	containerManager, err := container_manager.NewManager(config)
	if err != nil {
		log.Print("Container Manager Type not found")
	}

	agent := agent.NewAgent(containerManager, mqttClient)
	mqttClient.ConnectMQTTBroker(agent)

	// Register after connection
	agent.Register()

	defer mqttClient.CloseConnection()
	<-c
}
