package main

import (
	"Gin_test/ConnectDatabase"
	"Gin_test/rest"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	store, redis := ConnectDatabase.Connection()
	server := rest.Server{Store: store, Redis: redis}
	r := gin.Default()
	r.GET("/get-house/:id", server.GetHouseByID)
	r.GET("/get-all-houses", server.GetAllHouses)
	r.POST("/add-new-house", server.AddNewHouse)
	r.GET("/get-image/:id", server.GetImageByID)
	r.POST("/upload-image-to-redis", server.SaveToRedisFile)
	if err := r.Run("localhost:8080"); err != nil {
		return
	}
}
