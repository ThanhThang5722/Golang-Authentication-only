package main

import (
	"authentication/pkg/auth"
	"authentication/pkg/database"
	"authentication/routes"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	//DATABASE CONNECTION
	db := database.GetMongoInstance()
	defer db.Client.Disconnect(context.Background())
	fmt.Println("MONGODB CONNECTED")

	// GENERATE JWT SECRET KEY
	auth.GenerateJWTKey()

	router := gin.Default()
	//router.Use(middleware.CorsMiddleware)
	api := router.Group("/api")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Ping successful",
			})
		})
	}
	//ROUTER DEFINE
	routes.UserRouter(api)
	router.Run()
}
