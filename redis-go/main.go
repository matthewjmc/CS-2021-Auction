package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/idoko/rediboard/db"
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
)

func main() {
	database, err := db.NewDatabase(RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}
	router := initRouter(database)
	router.Run(ListenAddr)
}

func initRouter(database *db.Database) *gin.Engine {

	return r
}
