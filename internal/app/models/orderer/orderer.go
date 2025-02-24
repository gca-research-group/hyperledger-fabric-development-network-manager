package orderer

import (
	"errors"

	"gorm.io/gorm"
)

type Orderer struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Domain string
	Port   int
}

func (o *Orderer) FindAll(db *gorm.DB) ([]Orderer, error) {

	var orderers []Orderer

	err := db.Find(&orderers).Error

	return orderers, err
}

func (o *Orderer) FindById(db *gorm.DB, id uint) (Orderer, error) {
	var orderer Orderer

	if err := db.First(&orderer, id).Error; err != nil {
		return orderer, errors.New("RECORD_NOT_FOUND")
	}

	return orderer, nil
}

func (o *Orderer) Create(db *gorm.DB, orderer *Orderer) (*Orderer, error) {
	if orderer.Domain == "" {
		return nil, errors.New("DOMAIN_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	err := db.Create(&orderer).Error

	return orderer, err
}

func (o *Orderer) Update(db *gorm.DB, id uint, orderer *Orderer) (*Orderer, error) {

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

func (o *Orderer) Delete(db *gorm.DB, id uint) error {
	if _, err := o.FindById(db, id); err != nil {
		return err
	}

	db.Delete(&Orderer{}, id)

	return nil
}
