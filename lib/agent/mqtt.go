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
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/conf"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	Client   MQTT.Client
	Retained *bool
	Broker   config.BrokerConfig
}

func NewMQTTClient(brokerConfig config.BrokerConfig) *MQTTClient {
	return &MQTTClient{
		Broker: brokerConfig,
	}
}

func (client *MQTTClient) ConnectMQTTBroker(agent *Agent) {
	MQTT.DEBUG = log.New(os.Stdout, "", 0)

	hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://"+client.Broker.Host+":"+client.Broker.Port, "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	fmt.Println("Try to connect to: ", *server)
	topic := flag.String("topic", constants.TopicPrefix+conf.GetConf().Id, "Topic to subscribe to")
	qos := *flag.Int("qos", 2, "The QoS to subscribe to messages at")
	client.Retained = flag.Bool("retained", false, "Are the messages sent with the retained flag")
	fmt.Println(client.Retained)
	clientId := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientId).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(*topic, byte(qos), agent.onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client.Client = MQTT.NewClient(connOpts)
	if token := client.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *server)
	}
}

func (client *MQTTClient) CloseConnection() {
	client.Client.Disconnect(250)
	time.Sleep(1 * time.Second)
}

func (client *MQTTClient) Publish(topic string, message string, qos int) {
	client.Client.Publish(topic, byte(qos), *client.Retained, message)
}