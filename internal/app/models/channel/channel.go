package channel

import (
	"errors"
	"fmt"
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/peer"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type ChannelDto struct {
	ID             int    `form:"id"`
	Name           string `form:"name"`
	Peers          []int  `form:"peers"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}

type Channel struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `json:"name"`
	Peers     []*peer.Peer `json:"peers" gorm:"many2many:channel_peers;"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func (c *Channel) FindAll(db *gorm.DB, queryOptions sql.QueryOptions, queryParams ChannelDto) (http.Response[[]Channel], error) {

	var channels []Channel
	var total int64
	stmt := db.Model(&Channel{})

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

	response := http.Response[[]Channel]{}

	return *response.NewResponse(channels, queryOptions, int(total)), err
}

func (c *Channel) FindById(db *gorm.DB, id uint) (Channel, error) {
	var channel Channel

	if err := db.Preload("Peers").First(&channel, id).Error; err != nil {
		return channel, errors.New("RECORD_NOT_FOUND")
	}

	return channel, nil
}

func (c *Channel) Create(db *gorm.DB, channel *Channel) (*Channel, error) {
	if channel.Name == "" {
		return nil, errors.New("CHANNEL_NAME_CANNOT_BE_EMPTY")
	}

	err := db.Create(&channel).Error

	return channel, err
}

func (c *Channel) Update(db *gorm.DB, channel *Channel) (*Channel, error) {
	if channel.ID == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if channel.Name == "" {
		return nil, errors.New("CHANNEL_NAME_CANNOT_BE_EMPTY")
	}

	tx := db.Begin()

	if err := c.DeletePeers(tx, channel.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete peers: %w", err)
	}

	updatedChannel := Channel{ID: channel.ID}
	err := tx.Model(&updatedChannel).
		Where("id = ?", channel.ID).
		Updates(Channel{Name: channel.Name, UpdatedAt: time.Now().UTC()}).
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

func (c *Channel) Delete(db *gorm.DB, id uint) error {
	if _, err := c.FindById(db, id); err != nil {
		return err
	}

	if err := c.DeletePeers(db, id); err != nil {
		return err
	}

	err := db.Delete(&Channel{}, id).Error

	return err
}

func (c *Channel) DeletePeers(db *gorm.DB, id uint) error {
	var channel Channel

	if err := db.Preload("Peers").First(&channel, id).Error; err != nil {
		return err
	}

	if err := db.Model(&channel).Association("Peers").Clear(); err != nil {
		return err
	}

	return nil
}

func (c *ChannelDto) ToEntity() Channel {
	entity := Channel{}
	entity.ID = uint(c.ID)
	entity.Name = c.Name

	for _, id := range c.Peers {
		currentPeer := peer.Peer{}
		currentPeer.ID = uint(id)
		entity.Peers = append(entity.Peers, &currentPeer)
	}

	return entity
}
