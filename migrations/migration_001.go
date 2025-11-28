package main

import (
	"github.com/amirkhgraphic/go-arcaptcha-service/initializers"
	"github.com/amirkhgraphic/go-arcaptcha-service/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// AutoMigrate keeps the schema in sync with the User model.
	initializers.DB.AutoMigrate(&models.User{})
}
