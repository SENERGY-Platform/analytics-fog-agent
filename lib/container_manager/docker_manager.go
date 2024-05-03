package container_manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	mqtt "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	operatorEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"
)

type DockerManager struct {
	Broker             mqtt.BrokerConfig
	ContainerPullImage bool
	ContainerNetwork   string
}

func NewDockerManager(brokerHost string, brokerPort string, containerPullImage bool, containerNetwork string) *DockerManager {
	return &DockerManager{
		Broker: mqtt.BrokerConfig{
			Host: brokerHost, // must not be localhost, as container has no access 
			Port: brokerPort,
		},
		ContainerPullImage: containerPullImage,
		ContainerNetwork:   containerNetwork,
	}
}

func (manager *DockerManager) CreateAndStartOperator(ctx context.Context, operator operatorEntities.StartOperatorControlCommand) (containerId string, err error) {
	operatorConfig, inputTopics, config, err := StartOperatorConfigsToString(operator)
	if err != nil {
		return
	}
	env := []string{
		"INPUT=" + inputTopics,
		"CONFIG=" + config,
		"OPERATOR_CONFIG=" + operatorConfig,
		"BROKER_HOST=" + manager.Broker.Host,
		"BROKER_PORT=" + manager.Broker.Port,
	}

	containerId, err = manager.RunContainer(ctx, operator.ImageId, env, manager.ContainerPullImage, operator.Config.PipelineId, operator.Config.OperatorId)
	return
}

func (manager *DockerManager) RestartOperator(ctx context.Context, operatorID string) (err error) {
	return
}

func (manager *DockerManager) RemoveOperator(ctx context.Context, operatorId string) (err error) {
	return manager.RemoveContainer(ctx, operatorId)
}

func (manager *DockerManager) GetOperatorState(ctx context.Context, operatorID string) (state OperatorState, err error) {
	return
}

func (manager *DockerManager) UpdateOperator(ctx context.Context, operatorID string, updateRequest operatorEntities.StartOperatorControlCommand) (err error) {
	return
}
	
func (manager *DockerManager) GetOperatorStates(ctx context.Context) (states map[string]OperatorState, err error) {
	return
}

func (manager *DockerManager) PullImage(imageName string) {
	ctx := context.Background()
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}

func (manager *DockerManager) RunContainer(ctx context.Context, imageName string, env []string, pull bool, pipelineId string, operatorId string) (string, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return "", err
	}

	if pull == true {
		out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			if err.Error() == "repository name must be canonical" {
				out, err = cli.ImagePull(ctx, "docker.io/"+imageName, image.PullOptions{})
			}
			if err != nil {
				return "", err
			}
		}
		_, _ = io.Copy(os.Stdout, out)
	}
	network := manager.ContainerNetwork
	logging.Logger.Debugf("Try to create docker container %s", imageName)

	resp, err := cli.ContainerCreate(
		ctx, 
		&container.Config{
			Image: imageName,
			Env:   env,
		}, 
		&container.HostConfig{
			NetworkMode: container.NetworkMode(network),
			RestartPolicy: container.RestartPolicy{
				Name: "on-failure",
				MaximumRetryCount: 5,
			},
		}, 
		nil, 
		nil,
		"fog-"+pipelineId+"-"+operatorId,
	)
	if err != nil {
		logging.Logger.Debugf("Could not create container: %s", err.Error())
		return "", err
	}
	logging.Logger.Debugf("Try to start docker container with ID: %s", resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID, nil
}

func (manager *DockerManager) ListAllContainers() {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
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

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
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

	containers, err := cli.ContainerList(ctx, container.ListOptions{
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

	removeOptions := container.RemoveOptions{Force: true}

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

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
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

func (manager *DockerManager) RemoveContainer(ctx context.Context, id string) (err error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{
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

	removeOptions := container.RemoveOptions{Force: true}

	for _, ct := range containers {
		if id == ct.ID {
			fmt.Print("Remove container ", ct.ID[:10], "... ")
			if err := cli.ContainerRemove(ctx, ct.ID, removeOptions); err != nil {
				return err
			}
			fmt.Println("Success")
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Container: %s not found", id))
}
