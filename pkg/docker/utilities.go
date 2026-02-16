package docker

import (
	"context"
	"os"
	"os/exec"
	"strings"

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

func RemoveContainersInNetwork(network string) error {
	list := exec.Command(
		"docker",
		"network",
		"inspect",
		network,
		"--format",
		"{{ range .Containers }}{{ .Name }}{{ \"\\n\" }}{{ end }}",
	)

	out, err := list.Output()
	if err != nil {
		return err
	}

	containers := strings.Fields(string(out))

	if len(containers) == 0 {
		return nil
	}

	for _, c := range containers {
		rm := exec.Command("docker", "rm", "-f", c)
		rm.Stderr = os.Stderr

		if err := rm.Run(); err != nil {
			return err
		}
	}

	return nil
}
