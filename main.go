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
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/agent"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/relay"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/container_manager"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/mqtt"

	srv_base "github.com/SENERGY-Platform/go-service-base/srv-base"

	"github.com/joho/godotenv"
)

func main() {
	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	err := godotenv.Load()
	if err != nil {
		log.Print("Cant load .env file")
	}

	config, err := config.NewConfig("")
	if err != nil {
		log.Print("Error loading config")
	}

	logFile, err := logging.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *srv_base.LogFileError
		if errors.As(err, &logFileError) {
			ec = 1
			return
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	logging.Logger.Debugf("config: %s", srv_base.ToJsonStr(config))

	watchdog := srv_base.NewWatchdog(logging.Logger, syscall.SIGINT, syscall.SIGTERM)

	conf.InitConf(config.DataDir)

	mqttClient := mqtt.NewMQTTClient(config.Broker, logging.Logger)

	containerManager, err := container_manager.NewManager(config)
	if err != nil {
		logging.Logger.Debug("Container Manager Type not found")
	}

	agent := agent.NewAgent(containerManager, mqttClient, conf.GetConf())
	relayController := relay.NewRelayController(agent)

	mqttClient.ConnectMQTTBroker(relayController)

	go agent.Register()

	watchdog.RegisterStopFunc(func() error {
		mqttClient.CloseConnection()
		return nil
	})

	watchdog.Start()

	ec = watchdog.Join()

}
