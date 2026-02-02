package docker

import (
	"context"

	"github.com/docker/docker/client"
)

func IsDockerRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}

	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}
