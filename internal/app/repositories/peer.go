package repositories

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"gorm.io/gorm"
)

type PeerRepository struct {
	DB *gorm.DB
}

func NewPeerRepository(db *gorm.DB) *PeerRepository {
	return &PeerRepository{DB: db}
}

func (r *PeerRepository) Create(entity *models.Peer) error {
	return r.DB.Create(entity).Error
}

func (r *PeerRepository) Update(id int, updates map[string]interface{}) (models.Peer, error) {
	entity := models.Peer{}
	err := r.DB.Model(&entity).Where("id = ?", id).UpdateColumns(updates).Error

	return entity, err
}

func (r *PeerRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Peer{}, id).Error
}

func (r *PeerRepository) FindMany() ([]models.Peer, error) {
	var entitys []models.Peer
	err := r.DB.Find(&entitys).Error
	return entitys, err
}

func (r *PeerRepository) FindById(id uint) (*models.Peer, error) {
	var entity models.Peer
	err := r.DB.First(&entity, id).Error
	return &entity, err
}
