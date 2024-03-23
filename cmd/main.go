package main

import (
	"fmt"
	"main/config"
	handlers "main/pkg/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.SetUpApplication()

}

func main() {
	fmt.Println("Server starting!")
	router := gin.Default()
	router.GET("/*path", handlers.RedirectHandler)
	router.POST("/shorten", handlers.ShortenURLHandler)

	router.Run(":8081")

}
