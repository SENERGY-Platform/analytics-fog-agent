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
	"log"
	"os"
	"syscall"
	"time"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/agent"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/relay"
	"github.com/SENERGY-Platform/analytics-fog-agent/migrations"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/container_manager"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/storage"
	mqttLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"

	srv_base "github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/SENERGY-Platform/go-service-base/watchdog"

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
		ec = 1
		return
	}

	err = logging.InitLogger(os.Stdout, true)
	if err != nil {
		log.Printf("Error init logging: %s", err.Error())
		ec = 1
		return
	}
	
	logging.Logger.Debug("config: " + srv_base.ToJsonStr(config))

	watchdog := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	conf.InitConf(config.DataDir)

	logging.Logger.Debug("Create new database at " + config.Database.ConnectionURL)
	db, err := storage.NewDB(config.Database.ConnectionURL)
	if err != nil {
		logging.Logger.Error("Cant init DB", "error", err.Error())
		ec = 1
		return
	}
	defer db.Close()
	migrations.MigrateDb(config.Database.ConnectionURL)

	storageHandler := storage.New(db)

	mqttConfig := mqttLib.BrokerConfig(config.Broker)
	mqttClient := mqtt.NewMQTTClient(mqttConfig, logging.Logger)

	containerManager, err := container_manager.NewManager(config)
	if err != nil {
		logging.Logger.Error("Container Manager Type not found")
		ec = 1
		return
	}

	agent := agent.NewAgent(containerManager, mqttClient, conf.GetConf(), time.Duration(config.ControlOperatorTimeout), storageHandler)
	subscriptionHandler := relay.NewRelayController(agent)
	mqttClient.SetSubscriptionHandler(subscriptionHandler)

	mqttClient.ConnectMQTTBroker(nil, nil)
	
	go agent.Register()

	watchdog.RegisterStopFunc(func() error {
		mqttClient.CloseConnection()
		return nil
	})

	logging.Logger.Info("Agent is ready")

	watchdog.Start()

	ec = watchdog.Join()

}
