package application

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
)

type fakeExecutor struct {
	commands []string
	output   []byte
	failAt   int
}

func (f *fakeExecutor) ExecCommand(name string, args ...string) error {
	f.commands = append(f.commands, strings.Join(append([]string{name}, args...), " "))
	if f.failAt == len(f.commands) {
		return errors.New("injected command failure")
	}
	return nil
}

func (f *fakeExecutor) OutputCommand(name string, args ...string) ([]byte, error) {
	f.commands = append(f.commands, strings.Join(append([]string{name}, args...), " "))
	if f.failAt == len(f.commands) {
		return nil, errors.New("injected command failure")
	}
	return f.output, nil
}

func TestGenerateArtifactsRejectsNonEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "existing"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := NewWorkflows(&fakeExecutor{}).GenerateArtifacts(&config.Config{Output: dir})
	if err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("expected non-empty directory error, got %v", err)
	}
}

func TestStopNetworkUsesInjectedExecutor(t *testing.T) {
	exec := &fakeExecutor{output: []byte("peer0.example.org\norderer.example.org\n")}
	err := NewWorkflows(exec).StopNetwork(&config.Config{Network: "example"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"docker network inspect example --format {{ range .Containers }}{{ .Name }}{{ \"\\n\" }}{{ end }}",
		"docker rm -f peer0.example.org",
		"docker rm -f orderer.example.org",
	}
	if !reflect.DeepEqual(exec.commands, want) {
		t.Fatalf("commands mismatch\nwant: %#v\n got: %#v", want, exec.commands)
	}
}

func TestStartNetworkStopsAtFirstCommandFailure(t *testing.T) {
	exec := &fakeExecutor{failAt: 2}
	cfg := config.Config{
		Output: "artifacts",
		Organizations: []config.Organization{{
			Name: "Org1", Domain: "org1.example.org",
			Orderers: []config.Orderer{{Name: "Orderer", Subdomain: "orderer"}},
			Peers:    []config.Peer{{Name: "Peer", Subdomain: "peer0"}},
		}},
	}
	err := NewWorkflows(exec).StartNetwork(&cfg)
	if err == nil || !strings.Contains(err.Error(), "Start Orderers") {
		t.Fatalf("expected contextual orderer failure, got %v", err)
	}
	if len(exec.commands) != 2 {
		t.Fatalf("expected workflow to stop after two commands, got %d", len(exec.commands))
	}
}

func TestGenerateArtifactsMatchesGoldenFiles(t *testing.T) {
	output := t.TempDir()
	cfg := artifactTestConfig(output)
	if err := NewWorkflows(&fakeExecutor{}).GenerateArtifacts(&cfg); err != nil {
		t.Fatal(err)
	}

	files := []string{
		"configtx.yml",
		"network.yml",
		filepath.Join("org1.example.org", "ca.org1.example.org.yml"),
		filepath.Join("org1.example.org", "peer0.org1.example.org.yml"),
	}
	for _, name := range files {
		got, err := os.ReadFile(filepath.Join(output, name))
		if err != nil {
			t.Fatal(err)
		}
		golden := filepath.Join("testdata", "artifacts", name)
		if os.Getenv("UPDATE_GOLDEN") == "1" {
			if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(golden, got, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		want, err := os.ReadFile(golden)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("generated artifact differs from golden file: %s", name)
		}
	}
}

func artifactTestConfig(output string) config.Config {
	return config.Config{
		Output: output, Network: "example",
		Capabilities: config.Capabilities{Channel: "V2_0", Orderer: "V2_0", Application: "V2_0"},
		Organizations: []config.Organization{{
			Name: "Org1", Domain: "org1.example.org", Bootstrap: true,
			CertificateAuthority: config.CertificateAuthority{ExposePort: 7054, Version: "latest"},
			Orderers:             []config.Orderer{{Name: "Orderer", Subdomain: "orderer", Port: 7050, ExposePort: 7050, Version: "2.5.15"}},
			Peers:                []config.Peer{{Name: "Peer0", Subdomain: "peer0", Port: 7051, ExposePort: 7051, Version: "2.5.15", IsAnchor: true}},
		}},
		Profiles: []config.Profile{{Name: "Default", Organizations: []string{"Org1"}, Consensus: config.Consensus{Type: "etcdraft"}}},
		Channels: []config.Channel{{Name: "defaultchannel", Profile: "Default"}},
	}
}
