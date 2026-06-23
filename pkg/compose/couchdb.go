package compose

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

type CouchDBNode struct {
	*yaml.Node
}

func NewCouchDB(
	domain string,
	peerSubdomain string,
	network string,
) *CouchDBNode {
	node := yaml.MappingNode(
		yaml.ScalarNode("container_name"),
		yaml.ScalarNode(ResolveCouchDBDomain(peerSubdomain, domain)),
		yaml.ScalarNode("image"),
		yaml.ScalarNode(ResolveCouchDBImage()),
		yaml.ScalarNode("environment"),
		yaml.MappingNode(
			yaml.ScalarNode("COUCHDB_USER"),
			yaml.ScalarNode("admin"),
			yaml.ScalarNode("COUCHDB_PASSWORD"),
			yaml.ScalarNode("adminpw"),
		),
		yaml.ScalarNode("volumes"),
		yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("./%s/data/peers/%s/couchdb:/opt/couchdb/data", domain, peerSubdomain))),
		yaml.ScalarNode("networks"),
		yaml.SequenceNode(yaml.ScalarNode(network)),
		yaml.ScalarNode("healthcheck"),
		yaml.MappingNode(
			yaml.ScalarNode("test"),
			yaml.SequenceNode(
				yaml.ScalarNode("CMD"),
				yaml.ScalarNode("curl"),
				yaml.ScalarNode("-f"),
				yaml.ScalarNode("http://admin:adminpw@localhost:5984/"),
			).WithFlowStyle(),
			yaml.ScalarNode("interval"), yaml.ScalarNode("5s"),
			yaml.ScalarNode("timeout"), yaml.ScalarNode("5s"),
			yaml.ScalarNode("retries"), yaml.ScalarNode("10"),
		),
	)

	return &CouchDBNode{node}
}

func (pn *CouchDBNode) Build() *yaml.Node {
	return pn.Node
}
