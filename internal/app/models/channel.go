package models

import (
	"time"
)

type Channel struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Peers     []*Peer    `json:"peers" gorm:"many2many:channel_peers;"`
	Orderers  []*Orderer `json:"orderers" gorm:"many2many:channel_orderers;"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
