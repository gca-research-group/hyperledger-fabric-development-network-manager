package peer

import (
	"errors"
	"time"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/http"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type PeerDto struct {
	ID             int    `form:"id"`
	Name           string `form:"name"`
	Domain         string `form:"domain"`
	Port           int    `form:"port"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}

type Peer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Port      int       `json:"port"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (o *Peer) FindAll(db *gorm.DB, queryOptions sql.QueryOptions, queryParams PeerDto) (http.Response[[]Peer], error) {

	var peers []Peer
	var total int64
	stmt := db.Model(&Peer{})

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
	desc := true

	if queryParams.OrderBy != "" {
		column = schema.NamingStrategy{}.ColumnName("", queryParams.OrderBy)
	}

	if queryParams.OrderDirection != "" {
		desc = queryParams.OrderDirection == "desc"
	}

	stmt.Order(clause.OrderByColumn{Column: clause.Column{Name: column}, Desc: desc})

	err := stmt.Offset(queryOptions.Offset).Limit(queryOptions.Limit).Find(&peers).Error
	stmt.Count(&total)

	response := http.Response[[]Peer]{}

	return *response.NewResponse(peers, queryOptions, int(total)), err
}

func (o *Peer) FindById(db *gorm.DB, id uint) (Peer, error) {
	var peer Peer

	if err := db.First(&peer, id).Error; err != nil {
		return peer, errors.New("RECORD_NOT_FOUND")
	}

	return peer, nil
}

func (o *Peer) Create(db *gorm.DB, peer *Peer) (*Peer, error) {
	if peer.Domain == "" {
		return nil, errors.New("DOMAIN_CANNOT_BE_EMPTY")
	}

	if peer.Name == "" {
		return nil, errors.New("PEER_NAME_CANNOT_BE_EMPTY")
	}

	err := db.Create(&peer).Error

	return peer, err
}

func (o *Peer) Update(db *gorm.DB, peer Peer) (*Peer, error) {

	if peer.ID == 0 {
		return nil, errors.New("ID_CANNOT_BE_EMPTY")
	}

	if peer.Domain == "" {
		return nil, errors.New("PEER_DOMAIN_CANNOT_BE_EMPTY")
	}

	if peer.Port == 0 {
		return nil, errors.New("PORT_CANNOT_BE_EMPTY")
	}

	if peer.Name == "" {
		return nil, errors.New("PEER_NAME_CANNOT_BE_EMPTY")
	}

	_peer := Peer{}
	err := db.Model(&_peer).Where("id = ?", peer.ID).UpdateColumns(Peer{Name: peer.Name, Domain: peer.Domain, Port: peer.Port, UpdatedAt: time.Now().UTC()}).Error

	return &_peer, err
}

func (o *Peer) Delete(db *gorm.DB, id uint) error {
	if _, err := o.FindById(db, id); err != nil {
		return err
	}

	err := db.Delete(&Peer{}, id).Error

	return err
}

func (o *Peer) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreatedAt = time.Now().UTC()
	o.UpdatedAt = time.Now().UTC()
	return
}

func (o *Peer) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now().UTC()
	return
}
