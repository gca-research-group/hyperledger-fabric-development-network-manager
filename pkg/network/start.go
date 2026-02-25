package network

import "fmt"

func (f *Network) Start() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Start Certificate Authorities", f.RunCAContainers},
		{"Start Orderers", f.RunOrdererContainers},
		{"Start Peers", f.RunPeerContainers},
		{"Start Tools", f.RunToolsContainers},
	}

	for _, step := range steps {
		fmt.Printf(">>> Step: %s\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}
