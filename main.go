package main

import (
	"github.com/amirkhgraphic/go-arcaptcha-service/controllers"
	"github.com/amirkhgraphic/go-arcaptcha-service/initializers"
	"github.com/gin-gonic/gin"
)

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

	// listens on 0.0.0.0:8080 by default
	router.Run()
}
