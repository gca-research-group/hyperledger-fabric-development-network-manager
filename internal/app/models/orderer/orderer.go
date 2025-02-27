package orderer

import (
	"errors"
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"gorm.io/gorm"
)

type OrdererDto struct {
	ID     int    `form:"id"`
	Name   string `form:"name"`
	Domain string `form:"domain"`
	Port   int    `form:"port"`
}

type Orderer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Port      int       `json:"port"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (o *Orderer) FindAll(db *gorm.DB, queryOptions sql.QueryOptions, queryParams OrdererDto) (http.Response[[]Orderer], error) {

	var orderers []Orderer
	var total int64
	stmt := db.Model(&Orderer{})

	if queryParams.Domain != "" {
		stmt.Where("domain ilike ?", "%"+queryParams.Domain+"%")
	}

	if queryParams.Name != "" {
		stmt.Where("name ilike ?", "%"+queryParams.Name+"%")
	}

	if queryParams.Port != 0 {
		stmt.Where("port = ?", queryParams.Port)
	}

	if queryParams.ID != 0 {
		stmt.Where("id = ?", queryParams.ID)
	}

	err := stmt.Offset(queryOptions.Offset).Limit(queryOptions.Limit).Find(&orderers).Error
	stmt.Count(&total)

	response := http.Response[[]Orderer]{}

	return *response.NewResponse(orderers, queryOptions, int(total)), err
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
