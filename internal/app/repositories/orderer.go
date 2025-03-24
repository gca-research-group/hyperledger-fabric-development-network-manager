package repositories

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"gorm.io/gorm"
)

type OrdererRepository struct {
	DB *gorm.DB
}

func NewOrdererRepository(db *gorm.DB) *OrdererRepository {
	return &OrdererRepository{DB: db}
}

func (r *OrdererRepository) Create(orderer *models.Orderer) error {
	return r.DB.Create(orderer).Error
}

func (r *OrdererRepository) Update(id int, updates map[string]interface{}) (models.Orderer, error) {
	entity := models.Orderer{}
	err := r.DB.Model(&entity).Where("id = ?", id).UpdateColumns(updates).Error

	return entity, err
}

func (r *OrdererRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Orderer{}, id).Error
}

func (r *OrdererRepository) FindMany() ([]models.Orderer, error) {
	var orderers []models.Orderer
	err := r.DB.Find(&orderers).Error
	return orderers, err
}

func (r *OrdererRepository) FindById(id uint) (*models.Orderer, error) {
	var orderer models.Orderer
	err := r.DB.First(&orderer, id).Error
	return &orderer, err
}
