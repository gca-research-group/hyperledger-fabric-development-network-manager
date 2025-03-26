package database

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/repositories"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/services"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	email := "admin@admin.com"

	if existingUser, _ := userService.FindByEmail(email); existingUser.ID != 0 {
		return
	}

	user := models.User{
		Name:    "Admin",
		Email:   email,
		IsSuper: true,
	}

	user.Password = "admin"

	userService.Create(&user)
}
