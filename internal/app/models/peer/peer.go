package peer

import (
	"errors"

	"gorm.io/gorm"
)

type Peer struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Domain string
	Port   int
}

func (p *Peer) FindAll(db *gorm.DB) ([]Peer, error) {

	var orderers []Peer

	err := db.Find(&orderers).Error

	return orderers, err
}

func (p *Peer) FindById(db *gorm.DB, id uint) (Peer, error) {
	var peer Peer

	if err := db.First(&peer, id).Error; err != nil {
		return peer, errors.New("RECORD_NOT_FOUND")
	}

	return peer, nil
}

func (p *Peer) Create(db *gorm.DB, peer *Peer) (*Peer, error) {
	if peer.Domain == "" {
		return nil, errors.New("DOMAIN_CANNOT_BE_EMPTY")
	}

	if peer.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	err := db.Create(&peer).Error

	return peer, err
}

func (p *Peer) Update(db *gorm.DB, id uint, peer *Peer) (*Peer, error) {

	if id == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if peer.Domain != "" {
		db.Model(&peer).Update("domain", peer.Domain)
	}

	if peer.Port == 0 {
		return nil, errors.New("PORT_CANNOT_BE_EMPTY")
	}

	if peer.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	peer.ID = id

	err := db.Save(&peer).Error

	return peer, err
}

func (p *Peer) Delete(db *gorm.DB, id uint) error {
	if _, err := p.FindById(db, id); err != nil {
		return err
	}

	db.Delete(&Peer{}, id)

	return nil
}
