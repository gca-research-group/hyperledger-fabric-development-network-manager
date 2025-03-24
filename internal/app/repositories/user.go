package repositories

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(entity *models.User) error {
	return r.DB.Create(entity).Error
}

func (r *UserRepository) Update(id int, updates map[string]interface{}) (models.User, error) {
	entity := models.User{}
	err := r.DB.Model(&entity).Where("id = ?", id).UpdateColumns(updates).Error

	return entity, err
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) FindMany() ([]models.User, error) {
	var entitys []models.User
	err := r.DB.Find(&entitys).Error
	return entitys, err
}

func (r *UserRepository) FindById(id uint) (models.User, error) {
	var entity models.User
	err := r.DB.First(&entity, id).Error
	return entity, err
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var entity models.User
	err := r.DB.Model(&models.User{}).Where("email = ?", email).Find(&entity).Error
	return entity, err
}
