package services

import (
	"errors"
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/dtos"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type OrdererService struct {
	repository *repositories.OrdererRepository
}

func NewOrdererService(repository *repositories.OrdererRepository) *OrdererService {
	return &OrdererService{repository: repository}
}

func (s *OrdererService) FindById(id uint) (*models.Orderer, error) {
	orderer, err := s.repository.FindById(id)

	if err != nil {
		return nil, errors.New("REGISTER_NOT_FOUND")
	}

	return orderer, err
}

func (s *OrdererService) FindAll(queryOptions sql.QueryOptions, queryParams dtos.OrdererDto) (http.Response[[]models.Orderer], error) {

	var orderers []models.Orderer
	var total int64
	stmt := s.repository.DB.Model(&models.Orderer{})

	if queryParams.Domain != "" {
		stmt.Where("domain ilike ?", "%"+queryParams.Domain+"%")
	}

	if queryParams.Name != "" {
		stmt.Where("name ilike ?", "%"+queryParams.Name+"%")
	}

	if queryParams.ID != 0 {
		stmt.Where("id = ?", queryParams.ID)
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

	err := stmt.Offset(queryOptions.Offset).Limit(queryOptions.Limit).Find(&orderers).Error
	stmt.Count(&total)

	response := http.Response[[]models.Orderer]{}

	return *response.NewResponse(orderers, queryOptions, int(total)), err
}

func (s *OrdererService) Create(orderer *models.Orderer) (models.Orderer, error) {
	if orderer.Domain == "" {
		return models.Orderer{}, errors.New("DOMAIN_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		return models.Orderer{}, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	err := s.repository.Create(orderer)

	return models.Orderer{}, err
}

func (s *OrdererService) Update(orderer models.Orderer) (models.Orderer, error) {

	if orderer.ID == 0 {
		return models.Orderer{}, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if orderer.Domain == "" {
		return models.Orderer{}, errors.New("ORDERER_DOMAIN_CANNOT_BE_EMPTY")
	}

	if orderer.Name == "" {
		return models.Orderer{}, errors.New("ORDERER_NAME_CANNOT_BE_EMPTY")
	}

	if _, err := s.FindById(orderer.ID); err != nil {
		return models.Orderer{}, err
	}

	data := map[string]interface{}{
		"Name":      orderer.Name,
		"Domain":    orderer.Domain,
		"UpdatedAt": time.Now().UTC(),
	}

	entity, err := s.repository.Update(int(orderer.ID), data)

	return entity, err
}

func (s *OrdererService) Delete(id uint) error {
	if _, err := s.repository.FindById(id); err != nil {
		return err
	}

	err := s.repository.Delete(id)

	return err
}
