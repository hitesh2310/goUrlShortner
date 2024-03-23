package database

import (
	"log"
	"main/pkg/constants"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var RedisClient *redis.Client

var ctx = context.Background()

func SetUpRedis() {
	log.Println("In setup redis  func")
	RedisClient = EstablishConn()
	log.Println("Redis COnnection is Established Succesfully", RedisClient)

}

func EstablishConn() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     constants.ApplicationConfig.Redis.Host + ":" + constants.ApplicationConfig.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func GetRedisConn() *redis.Client {

	if RedisClient != nil {
		return EstablishConn()
	} else {
		return RedisClient
	}

}
