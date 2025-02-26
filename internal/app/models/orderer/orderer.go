package orderer

import (
	"errors"
	"log"
	"math"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"gorm.io/gorm"
)

type Orderer struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Port   int    `json:"port"`
}

func (o *Orderer) FindAll(db *gorm.DB, query http.Query) (http.Response[[]Orderer], error) {

	var orderers []Orderer
	var total int64
	hasMore := true

	err := db.Offset(query.Offset).Limit(query.Limit).Find(&orderers).Error
	db.Model(&Orderer{}).Count(&total)
	log.Printf("Offset %d, Limit %d", query.Offset, query.Limit)
	if (query.Offset + query.Limit) >= int(total) {
		hasMore = false
	}

	if len(orderers) == 0 {
		hasMore = false
	}

	return http.Response[[]Orderer]{
		HasMore: hasMore,
		Total:   int(total),
		Page:    (query.Offset / query.Limit) + 1,
		Pages:   int(math.Ceil(float64(total) / float64(query.Limit))),
		Data:    orderers,
	}, err
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
