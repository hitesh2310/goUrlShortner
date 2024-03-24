package database

import (
	"encoding/json"
	"fmt"
	"log"
	"main/logs"
	"main/pkg/constants"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var RedisClient *redis.Client

var ctx = context.Background()

func SetUpRedis() {
	fmt.Println("In setup redis  func")
	RedisClient = EstablishConn()
	fmt.Println("Redis COnnection is Established Succesfully", RedisClient)

}

func EstablishConn() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     constants.ApplicationConfig.Redis.Host + ":" + constants.ApplicationConfig.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis: ", err)
		return nil
	}

	return rdb
}

func GetRedisConn() *redis.Client {

	if RedisClient != nil {
		return RedisClient
	} else {
		return EstablishConn()
	}

}

func HGet(key, field string) (string, error) {
	if RedisClient == nil {
		RedisClient = GetRedisConn()
	}
	val, err := RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		log.Printf("Error retrieving value from Redis hash: %v\n", err)
	}
	return val, err
}

func HSetShortLongMapping(key, field, value string) error {
	if RedisClient == nil {
		RedisClient = GetRedisConn()
	}
	err := RedisClient.HSet(ctx, key, field, value).Err()
	if err != nil {
		log.Printf("Error setting value in Redis hash: %v\n", err)
	}
	return err
}

func HGetAll(key string) (map[string]string, error) {
	logs.InfoLog("In HGETALL func")
	if RedisClient == nil {
		RedisClient = GetRedisConn()
	}

	result, err := RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		logs.ErrorLog("Error retrieving values from Redis hash: %v\n", err)
		return nil, err
	}

	fmt.Println("Result:", result)
	values := make(map[string]string)
	for k, v := range result {
		values[k] = v
	}

	fmt.Println("Values", values)
	return values, nil
}

func UpdateCache() {
	for {
		logs.InfoLog("CACHE UPDATE")
		result, err := HGetAll("url_mapping")
		if err != nil {
			//slack Alert
			logs.ErrorLog("Error to get url_mappping result %v", err)
		}
		currMap := make(map[string]interface{})
		for key := range result {
			logs.InfoLog("Key:: %v", key)
			logs.InfoLog("Value:: %v", result[key])
			err := json.Unmarshal([]byte(result[key]), &currMap)
			if err != nil {
				logs.InfoLog("Error to unmarshal redis value into map: %v", err)
			}

			epochTime, ok := currMap["epochTime"].(float64)
			if !ok {
				logs.ErrorLog("Error: epochTime field not found or not a number")

			}

			epochTimeInSeconds := int64(epochTime)
			timeFromEpoch := time.Unix(epochTimeInSeconds, 0)

			// Calculate the difference between the current time and the epoch time
			currentTime := time.Now()
			timeDiff := currentTime.Sub(timeFromEpoch)

			if timeDiff > 2*time.Hour {
				logs.InfoLog("The epoch time is older than 2 hours, key::%v", key)
				HDel(key)
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func HDel(key string) {

	if RedisClient == nil {
		RedisClient = GetRedisConn()
	}

	// Retrieve all fields and values from the Redis hash
	_, err := RedisClient.HDel(ctx, "url_mapping", key).Result()
	if err != nil {
		logs.ErrorLog("Error deleting the key from cache %v", err)
		return
	}

}
