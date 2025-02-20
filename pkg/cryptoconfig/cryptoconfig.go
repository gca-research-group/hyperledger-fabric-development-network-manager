package cryptoconfig

import (
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-network-manager/pkg"
)

type Spec struct {
	Hostname string `yaml:"Hostname"`
}

type OrdererOrg struct {
	Name   string `yaml:"Name"`
	Domain string `yaml:"Domain"`
	Specs  []Spec `yaml:"Specs"`
}

type Template struct {
	Count int `yaml:"Count"`
}

type Users struct {
	Count int
}

type PeerOrg struct {
	Name          string   `yaml:"Name"`
	Domain        string   `yaml:"Domain"`
	EnableNodeOUs bool     `yaml:"EnableNodeOUs"`
	Template      Template `yaml:"Template"`
	Users         Users    `yaml:"Users"`
}

type CryptoConfig struct {
	OrdererOrgs []OrdererOrg `yaml:"OrdererOrgs"`
	PeerOrgs    []PeerOrg    `yaml:"PeerOrgs"`
}

func Build(config pkg.Config) CryptoConfig {
	var _cryptoconfig CryptoConfig

	for _, orderer := range config.Orderers {
		_cryptoconfig.OrdererOrgs = append(_cryptoconfig.OrdererOrgs, OrdererOrg{
			Name:   orderer.Name,
			Domain: orderer.Domain,
			Specs: []Spec{
				{Hostname: strings.Split(orderer.Domain, ".")[0]},
			},
		})
	}

	for _, peer := range config.Peers {

		if peer.Peers < 1 {
			peer.Peers = 1
		}

		if peer.Users < 1 {
			peer.Users = 1
		}

		_cryptoconfig.PeerOrgs = append(_cryptoconfig.PeerOrgs, PeerOrg{
			Name:          peer.Name,
			Domain:        peer.Domain,
			EnableNodeOUs: true,
			Template: Template{
				Count: peer.Peers,
			},
			Users: Users{
				Count: peer.Users,
			},
		})

		anchorPeerPort := peer.Port

		if anchorPeerPort < 1 {
			anchorPeerPort = 7051
		}
	}

	return _cryptoconfig
}
