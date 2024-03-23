package handlers

import (
	"encoding/json"
	"fmt"
	"main/pkg/constants"
	"main/pkg/database"
	"main/pkg/models"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RedirectHandler(c *gin.Context) {

	fmt.Println("code check ", c.Request.URL.RequestURI())
	// c.Header("content-type", "html")
	shortUrl := c.Request.URL.RequestURI()
	var cacheEntry models.RedisEntry
	cacheString, err := database.HGet("url_mapping", shortUrl[1:])
	if err != nil {
		fmt.Println("Error getting url mapping from cache: ", err)
		//check DB
		// database.GetEntry(shortUrl)
	} else {
		json.Unmarshal([]byte(cacheString), &cacheEntry)
		fmt.Println("CACHE ENTRY::", cacheEntry)
	}

	c.Redirect(http.StatusSeeOther, cacheEntry.LongUrl)

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
				c.JSON(http.StatusOK, gin.H{"shortUrl": "localhost:8081/" + shortUrl})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	constants.GlobalMutex.Lock()
	constants.Counter++
	constants.GlobalMutex.Unlock()

	//Add to Cache
	var cacheEntry models.RedisEntry
	cacheEntry.LongUrl = req.Url
	cacheEntry.EpochTime = int(time.Now().Unix())

	cacheEntryByte, err := json.Marshal(cacheEntry)
	if err != nil {
		fmt.Println("Failed to marshal cache entry")
	}
	cacheEntryString := string(cacheEntryByte)
	database.HSetShortLongMapping("url_mapping", encodedString, cacheEntryString)

	c.JSON(http.StatusOK, gin.H{"shortUrl": "localhost:8081/" + encodedString})
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
