package model

type Config struct {
	ID     uint   `gorm:"primaryKey"`
	Config string `gorm:"type:text"`
}
