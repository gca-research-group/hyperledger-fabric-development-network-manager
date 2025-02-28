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

func (o *Orderer) Update(db *gorm.DB, orderer Orderer) (*Orderer, error) {

	if orderer.ID == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if orderer.Domain == "" {
		return nil, errors.New("ORDERER_DOMAIN_CANNOT_BE_EMPTY")
	}

	if orderer.Port == 0 {
		return nil, errors.New("PORT_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		return nil, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	_orderer := Orderer{}
	err := db.Model(&_orderer).Where("id = ?", orderer.ID).UpdateColumns(Orderer{Name: orderer.Name, Domain: orderer.Domain, Port: orderer.Port, UpdatedAt: time.Now().UTC()}).Error

	return &_orderer, err
}

func (o *Orderer) Delete(db *gorm.DB, id uint) error {
	if _, err := o.FindById(db, id); err != nil {
		return err
	}

	err := db.Delete(&Orderer{}, id).Error

	return err
}

func (o *Orderer) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreatedAt = time.Now().UTC()
	o.UpdatedAt = time.Now().UTC()
	return
}

func (o *Orderer) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now().UTC()
	return
}
