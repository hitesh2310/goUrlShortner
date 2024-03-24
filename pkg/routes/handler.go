package handlers

import (
	"encoding/json"
	"main/logs"
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

	logs.InfoLog("URI::%v", c.Request.URL.RequestURI())

	shortUrl := c.Request.URL.RequestURI()
	var cacheEntry models.RedisEntry
	cacheString, err := database.HGet("url_mapping", shortUrl[1:])
	if err != nil {
		logs.ErrorLog("Error getting url mapping from cache: %v", err)
		//check DB
		longUrl := database.GetEntry(shortUrl)
		if longUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}
		cacheEntry.LongUrl = longUrl
		cacheEntry.EpochTime = int(time.Now().Unix())
		//adding to redis cache
		cacheEntryByte, err := json.Marshal(cacheEntry)
		if err != nil {
			logs.ErrorLog("Failed to marshal cache entry")
		}
		cacheEntryString := string(cacheEntryByte)
		database.HSetShortLongMapping("url_mapping", shortUrl, cacheEntryString)

	} else {
		json.Unmarshal([]byte(cacheString), &cacheEntry)
		logs.InfoLog("CACHE ENTRY::", cacheEntry)
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

	logs.InfoLog("Request recieved:: %v", req)
	logs.InfoLog("Curent counter ::%v", constants.Counter)
	encodedString := EncodeBase62(constants.Counter)

	logs.InfoLog("Encoding string::%v", encodedString)
	logs.InfoLog("COUNTER:: %v", constants.Counter)

	//insert into DB
	err := database.AddEntry(req.Url, encodedString)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			shortUrl, err := database.GetShortUrl(req.Url)
			if shortUrl == "" || err != nil {
				logs.ErrorLog("Error getting short url from database")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting short url"})
				return
			} else {
				logs.InfoLog("Got short url from database %v", shortUrl)
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
		logs.ErrorLog("Failed to marshal cache entry")
	}
	cacheEntryString := string(cacheEntryByte)
	database.HSetShortLongMapping("url_mapping", encodedString, cacheEntryString)
	database.IncrementHost(req.Url)
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

func StatHandler(c *gin.Context) {

	result := database.GetTop3()
	c.JSON(http.StatusOK, gin.H{"result": result})
}
