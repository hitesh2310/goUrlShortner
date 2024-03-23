package handlers

import (
	"fmt"
	"main/pkg/constants"
	"main/pkg/database"
	"main/pkg/models"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func RedirectHandler(c *gin.Context) {

	fmt.Println("code check ", c.Request.URL.RequestURI())
	// c.Header("content-type", "html")
	c.Redirect(http.StatusSeeOther, "http://www.google.com")

}

func ShortenURLHandler(c *gin.Context) {

	var req models.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidURL(req.Url) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return

	}

	if len(req.Url) >= 700 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL too long"})
		return
	}

	fmt.Println("Request recieved: ", req)
	fmt.Println("Curent counter ::", constants.Counter)
	encodedString := EncodeBase62(constants.Counter)

	fmt.Println("Encoding string::", encodedString)
	fmt.Println("COUNTER::", constants.Counter)

	//insert into DB
	err := database.AddEntry(req.Url, encodedString)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			shortUrl, err := database.GetShortUrl(req.Url)
			if shortUrl == "" || err != nil {
				fmt.Println("Error getting short url from database")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting short url"})
				return
			} else {
				fmt.Println("Got short url from database")
				c.JSON(http.StatusOK, gin.H{"shortUrl": "localhost:8081/short/" + shortUrl})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	constants.GlobalMutex.Lock()
	constants.Counter++
	constants.GlobalMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{"shortUrl": "localhost:8081/short/" + encodedString})
}

func isValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func EncodeBase62(num int) string {
	var result string
	base := len(constants.Charset)
	for num > 0 {
		remainder := num % base
		result = string(constants.Charset[remainder]) + result
		num = num / base
	}
	return result
}
