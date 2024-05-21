package config

import (
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	srv_base "github.com/SENERGY-Platform/go-service-base/srv-base"
	"github.com/y-du/go-log-level/level"
)


type Config struct {
	ContainerNetwork    string `json:"container_network" env_var:"CONTAINER_NETWORK"`
	ContainerBrokerHost string `json:"container_broker_host" env_var:"CONTAINER_BROKER_HOST"`
	ContainerBrokerPort string `json:"container_broker_port" env_var:"CONTAINER_BROKER_PORT"`
	Broker              mqtt.FogBrokerConfig
	ModuleManagerURL       string `json:"module_manager_url" env_var:"MODULE_MANAGER_URL"`
	ContainerPullImage  bool                  `json:"container_pull_image" env_var:"CONTAINER_PULL_IMAGE"`
	ContainerManager    string                `json:"container_manager" env_var:"CONTAINER_MANAGER"`
	Logger              srv_base.LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
	DataDir             string                `json:"data_dir" env_var:"DATA_DIR"`
	DeploymentID string `json:"deployment_id" env_var:"MGW_DID"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		ContainerNetwork:    "bridge",
		ContainerBrokerHost: "localhost",
		ContainerManager: constants.MGWManager,
		Broker: mqtt.FogBrokerConfig{
			Port: "1883",
			Host: "localhost",
		},
		ModuleManagerURL: "http://localhost",
		ContainerPullImage: true,
		Logger: srv_base.LoggerConfig{
			Level:        level.Debug,
			Utc:          true,
			Microseconds: true,
			Terminal:     true,
		},
		DataDir: "./data",
	}

	err := srv_base.LoadConfig(path, &cfg, nil, nil, nil)
	return &cfg, err
}
