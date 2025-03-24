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

type PeerService struct {
	Repository *repositories.PeerRepository
}

func NewPeerService(repository *repositories.PeerRepository) *PeerService {
	return &PeerService{Repository: repository}
}

func (s *PeerService) FindAll(queryOptions sql.QueryOptions, queryParams dtos.PeerDto) (http.Response[[]models.Peer], error) {

	var peers []models.Peer
	var total int64
	stmt := s.Repository.DB.Model(&models.Peer{})

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

	column := "name"
	desc := false

	if queryParams.OrderBy != "" {
		column = schema.NamingStrategy{}.ColumnName("", queryParams.OrderBy)
	}

	if queryParams.OrderDirection != "" {
		desc = queryParams.OrderDirection == "desc"
	}

	stmt.Order(clause.OrderByColumn{Column: clause.Column{Name: column}, Desc: desc})

	err := stmt.Offset(queryOptions.Offset).Limit(queryOptions.Limit).Find(&peers).Error
	stmt.Count(&total)

	response := http.Response[[]models.Peer]{}

	return *response.NewResponse(peers, queryOptions, int(total)), err
}

func (s *PeerService) FindById(id uint) (models.Peer, error) {
	entity, err := s.Repository.FindById(id)

	if err != nil {
		return *entity, errors.New("RECORD_NOT_FOUND")
	}

	return *entity, nil
}

func (s *PeerService) Create(peer *models.Peer) (*models.Peer, error) {
	if peer.Domain == "" {
		return nil, errors.New("DOMAIN_CANNOT_BE_EMPTY")
	}

	if peer.Name == "" {
		return nil, errors.New("PEER_NAME_CANNOT_BE_EMPTY")
	}

	err := s.Repository.Create(peer)

	return peer, err
}

func (s *PeerService) Update(peer models.Peer) (models.Peer, error) {

	if peer.ID == 0 {
		return models.Peer{}, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if peer.Domain == "" {
		return models.Peer{}, errors.New("PEER_DOMAIN_CANNOT_BE_EMPTY")
	}

	if peer.Port == 0 {
		return models.Peer{}, errors.New("PORT_CANNOT_BE_EMPTY")
	}

	if peer.Name == "" {
		return models.Peer{}, errors.New("PEER_NAME_CANNOT_BE_EMPTY")
	}

	if _, err := s.FindById(peer.ID); err != nil {
		return models.Peer{}, err
	}

	data := map[string]interface{}{
		"Name":      peer.Name,
		"Domain":    peer.Domain,
		"Port":      peer.Port,
		"UpdatedAt": time.Now().UTC(),
	}

	entity, err := s.Repository.Update(int(peer.ID), data)

	return entity, err
}

func (s *PeerService) Delete(id uint) error {
	if _, err := s.FindById(id); err != nil {
		return err
	}

	err := s.Repository.Delete(id)

	return err
}
