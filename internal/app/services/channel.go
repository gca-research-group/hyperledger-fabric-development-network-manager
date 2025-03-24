package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/dtos"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type ChannelService struct {
	Repository *repositories.ChannelRepository
}

func NewChannelService(repository *repositories.ChannelRepository) *ChannelService {
	return &ChannelService{Repository: repository}
}

func (s *ChannelService) FindAll(queryOptions sql.QueryOptions, queryParams dtos.ChannelDto) (http.Response[[]models.Channel], error) {

	var channels []models.Channel
	var total int64
	stmt := s.Repository.DB.Model(&models.Channel{})

	if queryParams.Name != "" {
		stmt.Where("name ilike ?", "%"+queryParams.Name+"%")
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

	err := stmt.Offset(queryOptions.Offset).Limit(queryOptions.Limit).Find(&channels).Error
	stmt.Count(&total)

	response := http.Response[[]models.Channel]{}

	return *response.NewResponse(channels, queryOptions, int(total)), err
}

func (s *ChannelService) FindById(id uint) (models.Channel, error) {
	var channel models.Channel

	if err := s.Repository.DB.Preload("Peers").First(&channel, id).Error; err != nil {
		return channel, errors.New("RECORD_NOT_FOUND")
	}

	return channel, nil
}

func (s *ChannelService) Create(channel *models.Channel) (*models.Channel, error) {
	if channel.Name == "" {
		return nil, errors.New("CHANNEL_NAME_CANNOT_BE_EMPTY")
	}

	err := s.Repository.Create(channel)

	return channel, err
}

func (s *ChannelService) Update(channel *models.Channel) (*models.Channel, error) {
	if channel.ID == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if channel.Name == "" {
		return nil, errors.New("CHANNEL_NAME_CANNOT_BE_EMPTY")
	}

	tx := s.Repository.DB.Begin()

	if err := s.DeletePeers(channel.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete peers: %w", err)
	}

	updatedChannel := models.Channel{ID: channel.ID}
	err := tx.Model(&updatedChannel).
		Where("id = ?", channel.ID).
		Updates(models.Channel{Name: channel.Name, UpdatedAt: time.Now().UTC()}).
		Error

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	if len(channel.Peers) > 0 {
		if err := tx.Model(&updatedChannel).Association("Peers").Append(channel.Peers); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update peers: %w", err)
		}
	}

	if err := tx.Preload("Peers").First(&updatedChannel, channel.ID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to reload updated channel: %w", err)
	}

	tx.Commit()
	return &updatedChannel, nil
}

func (s *ChannelService) Delete(id uint) error {
	if _, err := s.Repository.FindById(id); err != nil {
		return err
	}

	if err := s.DeletePeers(id); err != nil {
		return err
	}

	err := s.Repository.Delete(id)

	return err
}

func (s *ChannelService) DeletePeers(id uint) error {
	var channel models.Channel

	if err := s.Repository.DB.Preload("Peers").First(&channel, id).Error; err != nil {
		return err
	}

	if err := s.Repository.DB.Model(&channel).Association("Peers").Clear(); err != nil {
		return err
	}

	return nil
}
