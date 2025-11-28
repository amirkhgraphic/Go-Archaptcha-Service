package main

import (
	"github.com/amirkhgraphic/go-arcaptcha-service/controllers"
	"github.com/amirkhgraphic/go-arcaptcha-service/initializers"
	_ "github.com/amirkhgraphic/go-arcaptcha-service/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Arcaptcha Service API
// @version 1.0
// @BasePath /

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	// Health/ping endpoint kept for quick checks.
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Fake arcaptcha helpers for local testing.
	fake := router.Group("/__fake")
	{
		fake.GET("/arcaptcha/challenge", controllers.GenerateFakeChallenge)
		fake.POST("/arcaptcha/verify", controllers.VerifyFakeChallenge)
	}

	api := router.Group("/api")
	{
		api.POST("/users", controllers.CreateUser)
		api.GET("/users", controllers.ListUsers)
		api.GET("/users/:id", controllers.GetUser)
		api.PATCH("/users/:id", controllers.UpdateUser)

		api.GET("/users/group", controllers.GroupUsers)
	}

	// Serve swagger UI (uses the bundled docs/swagger.json)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// listens on 0.0.0.0:8080 by default
	router.Run()
}
