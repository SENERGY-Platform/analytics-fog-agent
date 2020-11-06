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
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

func PullImage(imageName string) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
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

func RunContainer(imageName string, env []string, pull bool, pipelineId string, operatorId string) (string, error) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		return "", err
	}

	if pull == true {
		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		_, _ = io.Copy(os.Stdout, out)
	}
	network := GetEnv("CONTAINER_NETWORK", "bridge")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
	}, &container.HostConfig{NetworkMode: container.NetworkMode(network)}, nil, "fog-"+pipelineId+"-"+operatorId)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID, nil
}

func ListAllContainers() {
	cli, err := docker.NewEnvClient()
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

func StopAllContainers() {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ct := range containers {
		fmt.Print("Stopping container ", ct.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, ct.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}

func RemoveAllContainers() {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   false,
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

func StopContainer(id string) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
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
			if err := cli.ContainerStop(ctx, ct.ID, nil); err != nil {
				panic(err)
			}
			fmt.Println("Success")
		}
	}
}

func RemoveContainer(id string) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   false,
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
		if id == ct.ID {
			fmt.Print("Remove container ", ct.ID[:10], "... ")
			if err := cli.ContainerRemove(ctx, ct.ID, removeOptions); err != nil {
				panic(err)
			}
			fmt.Println("Success")
		}
	}
}
