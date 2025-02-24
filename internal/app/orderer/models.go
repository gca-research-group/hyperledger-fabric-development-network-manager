package orderer

import (
	"log"

	"gorm.io/gorm"
)

type Orderer struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Domain string
	Port   int
}

func (o *Orderer) FindAll(db *gorm.DB) []Orderer {

	var orderers []Orderer

	if err := db.Find(&orderers).Error; err != nil {
		log.Fatal("Error selecting records", err)
	}

	return orderers
}

func (o *Orderer) FindById(db *gorm.DB, id int) Orderer {
	var orderer Orderer

	if err := db.First(&orderer, id).Error; err != nil {
		log.Fatal("Error selecting record", err)
	}

	return orderer
}

func (o *Orderer) CreateOrUpdate(db *gorm.DB, orderer *Orderer) *Orderer {
	if orderer.Domain == "" {
		log.Fatal("DOMAIN_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		log.Fatal("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	if err := db.Save(&orderer).Error; err != nil {
		log.Fatal("ERROR_CREATING_RECORD", err)
	}

	return orderer
}

func (o *Orderer) Delete(db *gorm.DB, id int) {
	db.Delete(&Orderer{}, id)
}
