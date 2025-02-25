package channel

import (
	"errors"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/peer"
	"gorm.io/gorm"
)

type Channel struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Peers []*peer.Peer `gorm:"many2many:channel_peers;"`
}

func (c *Channel) FindAll(db *gorm.DB) ([]Channel, error) {

	var orderers []Channel

	err := db.Preload("Peers").Find(&orderers).Error

	return orderers, err
}

func (c *Channel) FindById(db *gorm.DB, id uint) (Channel, error) {
	var channel Channel

	if err := db.Preload("Peers").First(&channel, id).Error; err != nil {
		return channel, errors.New("RECORD_NOT_FOUND")
	}

	return channel, nil
}

func (c *Channel) Create(db *gorm.DB, channel *Channel) (*Channel, error) {
	if channel.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	err := db.Create(&channel).Error

	return channel, err
}

func (c *Channel) Update(db *gorm.DB, id uint, channel *Channel) (*Channel, error) {

	if id == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if channel.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	c.DeletePeers(db, id)

	channel.ID = id

	if err := db.Save(&channel).Error; err != nil {
		return nil, err
	}

	db.Preload("Peers").First(&channel, id)

	return channel, nil
}

func (c *Channel) Delete(db *gorm.DB, id uint) error {
	c.DeletePeers(db, id)

	db.Delete(&Channel{}, id)

	return nil
}

func (c *Channel) DeletePeers(db *gorm.DB, id uint) error {
	var channel Channel

	if err := db.Preload("Peers").First(&channel, id).Error; err != nil {
		return err
	}

	if err := db.Model(&channel).Association("Peers").Clear(); err != nil {
		return err
	}

	return nil
}
