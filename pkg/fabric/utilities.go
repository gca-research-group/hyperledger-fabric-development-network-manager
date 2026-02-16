package fabric

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

func buildToolsContainerName(organization config.Organization) string {
	return fmt.Sprintf("hyperledger-fabric-tools-%s", strings.ToLower(organization.Name))
}
