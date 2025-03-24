package dtos

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"

type ChannelDto struct {
	ID             int    `form:"id"`
	Name           string `form:"name"`
	Peers          []int  `form:"peers"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}

func (d *ChannelDto) ToEntity() models.Channel {
	entity := models.Channel{}
	entity.ID = uint(d.ID)
	entity.Name = d.Name

	for _, id := range d.Peers {
		currentPeer := models.Peer{}
		currentPeer.ID = uint(id)
		entity.Peers = append(entity.Peers, &currentPeer)
	}

	return entity
}
