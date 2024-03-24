package main

import (
	"fmt"
	"main/config"
	"main/pkg/database"
	handlers "main/pkg/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.SetUpApplication()

}

func main() {
	go func() {
		database.UpdateCache()
	}()
	fmt.Println("Server starting!")
	router := gin.Default()
	router.GET("/stat", handlers.StatHandler)
	router.POST("/shorten", handlers.ShortenURLHandler)
	router.GET("/:path", handlers.RedirectHandler)
	router.Run(":8081")

}
