package repositories

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	DB *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{DB: db}
}

func (r *ChannelRepository) Create(entity *models.Channel) error {
	return r.DB.Create(entity).Error
}

func (r *ChannelRepository) Update(id int, updates map[string]interface{}) (models.Channel, error) {
	entity := models.Channel{}
	err := r.DB.Model(&entity).Where("id = ?", id).UpdateColumns(updates).Error

	return entity, err
}

func (r *ChannelRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Channel{}, id).Error
}

func (r *ChannelRepository) FindMany() ([]models.Channel, error) {
	var orderers []models.Channel
	err := r.DB.Find(&orderers).Error
	return orderers, err
}

func (r *ChannelRepository) FindById(id uint) (*models.Channel, error) {
	var entity models.Channel
	err := r.DB.First(&entity, id).Error
	return &entity, err
}
