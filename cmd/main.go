package main

import (
	"fmt"
	handlers "main/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Server starting!")

	router := gin.Default()

	router.POST("/shorten", handlers.ShortenURLHandler)

	router.GET("/short", handlers.RedirectHandler)
	router.Run(":8081")

}