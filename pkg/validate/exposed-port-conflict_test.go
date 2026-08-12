package validate

import "testing"

func TestExposedPortConflictFn(t *testing.T) {
	ports := make(map[int]portOwner)
	peer := portOwner{ownerType: "peer", name: "peer0"}
	assertNoError(t, ExposedPortConflictFn(7051, peer, ports))
	if got := ports[7051]; got != peer {
		t.Fatalf("expected port owner %#v, got %#v", peer, got)
	}
	assertValidationError(t, ExposedPortConflictFn(7051, portOwner{ownerType: "orderer", name: "orderer0"}, ports), RuleExposedPortConflict, "Exposed Port Conflict", "Port 7051 is assigned to both peer 'peer0' and orderer 'orderer0'.")
	assertNoError(t, ExposedPortConflictFn(0, portOwner{ownerType: "peer", name: "ignored"}, ports))
	if _, exists := ports[0]; exists {
		t.Fatal("zero port must not be recorded")
	}
}
