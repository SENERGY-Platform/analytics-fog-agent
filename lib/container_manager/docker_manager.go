package container_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/config"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

type DockerManager struct {
	Broker             config.BrokerConfig
	ContainerPullImage bool
	ContainerNetwork   string
}

func NewDockerManager(brokerHost string, brokerPort string, containerPullImage bool, containerNetwork string) *DockerManager {
	return &DockerManager{
		Broker: config.BrokerConfig{
			Host: brokerHost,
			Port: brokerPort,
		},
		ContainerPullImage: containerPullImage,
		ContainerNetwork:   containerNetwork,
	}
}

func (manager *DockerManager) StartOperator(operator operatorEntities.StartOperatorMessage) (containerId string, err error) {
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
		"BROKER_HOST=" + manager.Broker.Host,
		"BROKER_PORT=" + manager.Broker.Port,
	}

	containerId, err = manager.RunContainer(operator.ImageId, env, manager.ContainerPullImage, operator.Config.PipelineId, operator.Config.OperatorId)
	return
}

func (manager *DockerManager) StopOperator(operatorId string) (err error) {
	return manager.RemoveContainer(operatorId)
}

func (manager *DockerManager) PullImage(imageName string) {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}

func (manager *DockerManager) RunContainer(imageName string, env []string, pull bool, pipelineId string, operatorId string) (string, error) {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return "", err
	}

	if pull == true {
		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			if err.Error() == "repository name must be canonical" {
				out, err = cli.ImagePull(ctx, "docker.io/"+imageName, types.ImagePullOptions{})
			}
			if err != nil {
				return "", err
			}
		}
		_, _ = io.Copy(os.Stdout, out)
	}
	network := manager.ContainerNetwork
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
	}, &container.HostConfig{NetworkMode: container.NetworkMode(network)}, nil, nil,
		"fog-"+pipelineId+"-"+operatorId)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	//
	return resp.ID, nil
}

func (manager *DockerManager) ListAllContainers() {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ct := range containers {
		fmt.Println(ct.ID)
	}
}

func (manager *DockerManager) StopAllContainers() {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ct := range containers {
		fmt.Print("Stopping container ", ct.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, ct.ID, container.StopOptions{}); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}

func (manager *DockerManager) RemoveAllContainers() {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.Args{},
	})
	if err != nil {
		panic(err)
	}

	removeOptions := types.ContainerRemoveOptions{Force: true}

	for _, ct := range containers {
		fmt.Print("Remove container ", ct.ID[:10], "... ")
		if err := cli.ContainerRemove(ctx, ct.ID, removeOptions); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}

func (manager *DockerManager) StopContainer(id string) {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ct := range containers {
		if id == ct.ID {
			fmt.Print("Stopping container ", ct.ID[:10], "... ")
			if err := cli.ContainerStop(ctx, ct.ID, container.StopOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("Success")
		}
	}
}

func (manager *DockerManager) RemoveContainer(id string) (err error) {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.Args{},
	})
	if err != nil {
		return err
	}

	removeOptions := types.ContainerRemoveOptions{Force: true}

	for _, ct := range containers {
		if id == ct.ID {
			fmt.Print("Remove container ", ct.ID[:10], "... ")
			if err := cli.ContainerRemove(ctx, ct.ID, removeOptions); err != nil {
				return err
			}
			fmt.Println("Success")
		}
	}

	return nil
}
