package config

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	srv_base "github.com/SENERGY-Platform/go-service-base/srv-base"
)

type BrokerConfig struct {
	Host string `json:"broker_host" env_var:"BROKER_HOST"`
	Port string `json:"broker_port" env_var:"BROKER_PORT"`
}

type ModuleManagerConfig struct {
	Host string `json:"module_manager_host" env_var:"MODULE_MANAGER_HOST"`
	Port string `json:"module_manager_port" env_var:"MODULE_MANAGER_PORT"`
}

type Config struct {
	ContainerNetwork    string `json:"container_network" env_var:"CONTAINER_NETWORK"`
	ContainerBrokerHost string `json:"container_broker_host" env_var:"CONTAINER_BROKER_HOST"`
	Broker              BrokerConfig
	ModuleManager       ModuleManagerConfig
	ContainerPullImage  bool   `json:"container_pull_image" env_var:"CONTAINER_PULL_IMAGE"`
	ContainerManager    string `json:"container_manager" env_var:"CONTAINER_MANAGER"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		ContainerNetwork:    "bridge",
		ContainerBrokerHost: "localhost",
		Broker: BrokerConfig{
			Port: "1883",
			Host: "localhost",
		},
		ModuleManager: ModuleManagerConfig{
			Host: "localhost",
			Port: "8080",
		},
		ContainerPullImage: true,
		ContainerManager:   constants.DockerManager,
	}

	err := srv_base.LoadConfig(path, &cfg, nil, nil, nil)
	return &cfg, err
}
