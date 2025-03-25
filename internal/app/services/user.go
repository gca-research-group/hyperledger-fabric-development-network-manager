package services

import (
	"errors"
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/dtos"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

func (s *UserService) Create(entity *models.User) (*models.User, error) {
	if entity.Name == "" {
		return nil, errors.New("USER_NAME_CANNOT_BE_EMPTY")
	}

	if entity.Email == "" {
		return nil, errors.New("USER_EMAIL_CANNOT_BE_EMPTY")
	}

	if entity.Password == "" {
		return nil, errors.New("USER_PASSWORD_CANNOT_BE_EMPTY")
	}

	if existingUser, _ := s.Repository.FindByEmail(entity.Email); existingUser.ID != 0 {
		return nil, errors.New("USER_ALREADY_EXISTS")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)

	if err != nil {
		return entity, err
	}

	entity.Password = string(hashedPassword)

	err = s.Repository.Create(entity)

	return entity.Sanitize(), err
}

func (s *UserService) Delete(id uint) error {
	if _, err := s.Repository.FindById(id); err != nil {
		return err
	}

	err := s.Repository.Delete(id)

	return err
}

func (s *UserService) FindAll(queryOptions sql.QueryOptions, queryParams dtos.UserDto) (http.Response[[]models.User], error) {

	var users []models.User
	var total int64
	stmt := s.Repository.DB.Model(&models.User{})

	if queryParams.Name != "" {
		stmt.Where("name ilike ?", "%"+queryParams.Name+"%")
	}

	if queryParams.Email != "" {
		stmt.Where("email ilike ?", "%"+queryParams.Email+"%")
	}

	column := "name"
	desc := false

	if queryParams.OrderBy != "" {
		column = schema.NamingStrategy{}.ColumnName("", queryParams.OrderBy)
	}

	if queryParams.OrderDirection != "" {
		desc = queryParams.OrderDirection == "desc"
	}

	stmt.Order(clause.OrderByColumn{Column: clause.Column{Name: column}, Desc: desc})

	err := stmt.Offset(queryOptions.Offset).Limit(queryOptions.Limit).Find(&users).Error
	stmt.Count(&total)

	response := http.Response[[]models.User]{}

	return *response.NewResponse(users, queryOptions, int(total)), err
}

func (s *UserService) FindByEmail(email string) (models.User, error) {
	entity, err := s.Repository.FindByEmail(email)

	if err != nil || entity.ID == 0 {
		return entity, errors.New("RECORD_NOT_FOUND")
	}

	return entity, err
}

func (s *UserService) FindById(id uint) (*models.User, error) {
	entity, err := s.Repository.FindById(id)

	if err != nil {
		return &entity, errors.New("RECORD_NOT_FOUND")
	}

	return entity.Sanitize(), nil
}

func (s *UserService) Update(entity models.User) (*models.User, error) {
	if entity.ID == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if entity.Name == "" {
		return nil, errors.New("USER_NAME_CANNOT_BE_EMPTY")
	}

	if entity.Email == "" {
		return nil, errors.New("USER_EMAIL_CANNOT_BE_EMPTY")
	}

	if entity.Password == "" {
		return nil, errors.New("USER_PASSWORD_CANNOT_BE_EMPTY")
	}

	hashedPassword, err := entity.HashPassword(entity.Password)

	if err != nil {
		return &entity, err
	}

	_user := models.User{}
	err = s.Repository.DB.Model(&_user).Where("id = ?", entity.ID).UpdateColumns(models.User{
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  hashedPassword,
		IsSuper:   entity.IsSuper,
		UpdatedAt: time.Now().UTC()}).Error

	return _user.Sanitize(), err
}
