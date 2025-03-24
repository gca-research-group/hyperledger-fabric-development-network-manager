package models

import (
	"time"
)

type Peer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Port      int       `json:"port"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
