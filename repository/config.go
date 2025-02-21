package repository

import (
	"encoding/json"
	"log"

	"github.com/gca-research-group/hyperledger-fabric-network-manager/config/database"
	"github.com/gca-research-group/hyperledger-fabric-network-manager/model"
	"github.com/gca-research-group/hyperledger-fabric-network-manager/pkg"
)

func Create(config pkg.Config) model.Config {

	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Error serializing config: %v", err)
	}

	_converted := string(configBytes)

	_config := model.Config{Config: _converted}

	connection := database.Connection()

	result := connection.Create(&_config)

	if result.Error != nil {
		log.Fatalf("failed to create user: %v", result.Error)
	}

	return _config

}
