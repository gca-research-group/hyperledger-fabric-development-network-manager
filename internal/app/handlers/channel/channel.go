package channel

import (
	"net/http"
	"strconv"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/errors"
	model "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/channel"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/peer"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChannelDto struct {
	Name  string
	Peers []int
}

func (c *ChannelDto) toEntity() model.Channel {
	entity := model.Channel{}
	entity.Name = c.Name

	for _, id := range c.Peers {
		currentPeer := peer.Peer{}
		currentPeer.ID = uint(id)
		entity.Peers = append(entity.Peers, &currentPeer)
	}

	return entity
}

func Index(c *gin.Context, db *gorm.DB) {
	entity := model.Channel{}
	data, err := entity.FindAll(db)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Show(c *gin.Context, db *gorm.DB) {
	entity := model.Channel{}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := entity.FindById(db, uint(id))

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func Create(c *gin.Context, db *gorm.DB) {
	var data ChannelDto

	if err := c.ShouldBindJSON(&data); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	entity := model.Channel{}
	channel := data.toEntity()
	_, err := entity.Create(db, &channel)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

func Update(c *gin.Context, db *gorm.DB) {
	var data ChannelDto
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	entity := model.Channel{}
	channel := data.toEntity()
	_, err := entity.Update(db, uint(id), &channel)

	if err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, channel)
}

func Delete(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))

	entity := model.Channel{}
	if err := entity.Delete(db, uint(id)); err != nil {
		c.Error(&errors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
