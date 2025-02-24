package orderer

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

func (o *Peer) FindAll(db *gorm.DB) ([]Peer, error) {

	var orderers []Peer

	err := db.Find(&orderers).Error

	return orderers, err
}

func (o *Peer) FindById(db *gorm.DB, id uint) (Peer, error) {
	var orderer Peer

	if err := db.First(&orderer, id).Error; err != nil {
		return orderer, errors.New("RECORD_NOT_FOUND")
	}

	return orderer, nil
}

func (o *Peer) Create(db *gorm.DB, orderer *Peer) (*Peer, error) {
	if orderer.Domain == "" {
		return nil, errors.New("DOMAIN_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	err := db.Create(&orderer).Error

	return orderer, err
}

func (o *Peer) Update(db *gorm.DB, id uint, orderer *Peer) (*Peer, error) {

	if id == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if orderer.Domain != "" {
		db.Model(&orderer).Update("domain", orderer.Domain)
	}

	if orderer.Port == 0 {
		return nil, errors.New("PORT_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	orderer.ID = id

	err := db.Save(&orderer).Error

	return orderer, err
}

func (o *Peer) Delete(db *gorm.DB, id uint) error {
	if _, err := o.FindById(db, id); err != nil {
		return err
	}

	db.Delete(&Peer{}, id)

	return nil
}
