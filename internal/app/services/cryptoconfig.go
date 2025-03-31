package services

import (
	"log"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/cryptoconfig"
	"gopkg.in/yaml.v3"
)

type CryptoConfigService struct {
	channelService *ChannelService
}

func NewCryptoConfigService(channelService *ChannelService) *CryptoConfigService {
	return &CryptoConfigService{channelService: channelService}
}

func (s *CryptoConfigService) GenerateCryptoConfig(channelId int) ([]byte, error) {
	channel, err := s.channelService.FindById(uint(channelId))

	if err != nil {
		return nil, err
	}

	orderers := []pkg.Orderer{}
	peers := []pkg.Peer{}

	for _, orderer := range channel.Orderers {
		orderers = append(orderers, pkg.Orderer{
			Name:   orderer.Name,
			Domain: orderer.Domain,
			Port:   7050,
		})
	}

	for _, peer := range channel.Peers {
		peers = append(peers, pkg.Peer{
			Name:   peer.Name,
			Domain: peer.Domain,
			Port:   7050,
		})
	}

	config := cryptoconfig.Build(pkg.Config{
		Orderers: orderers,
		Peers:    peers,
	})

	out, err := yaml.Marshal(config)

	if err != nil {
		log.Fatal(err)
	}

	return out, nil
}
