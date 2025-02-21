package database

import (
	"log"

	"github.com/gca-research-group/hyperledger-fabric-network-manager/model"
)

func Migrate() {
	connection := Connection()

	err := connection.AutoMigrate(&model.Config{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
