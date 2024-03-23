package handlers

import (
	"fmt"
	"main/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RedirectHandler(c *gin.Context) {

	fmt.Println("", c.Request.URL.RequestURI())
	// c.Header("content-type", "html")
	c.Redirect(http.StatusSeeOther, "http://www.google.com")

}

func ShortenURLHandler(c *gin.Context) {

	var req models.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Request recieved: ", req)

}
