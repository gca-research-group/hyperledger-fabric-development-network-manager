package network

import (
	"fmt"
	"log/slog"
)

func (f *Network) Deploy() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Start Certificate Authorities", f.containerManager.RunCAContainers},
		{"Generate Certificates", f.identityManager.GenerateAll},

		{"Start Tools", f.containerManager.RunToolsContainers},

		{"Generate Genesis", f.GenerateGenesisBlock},

		{"Start Orderers", f.containerManager.RunOrdererContainers},
		{"Start Peers", f.containerManager.RunPeerContainers},
		{"Join Orderers", f.JoinOrdererToTheChannel},
		{"Fetch Genesis Block", f.FetchGenesisBlock},
		{"Join Peers", f.JoinPeersToTheChannels},
	}

	for _, step := range steps {
		slog.Info("Executing step", "step", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}
