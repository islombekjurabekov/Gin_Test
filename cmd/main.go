package main

import (
	"Gin_test/ConnectDatabase"
	"Gin_test/rest"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	store := ConnectDatabase.Connection()
	server := rest.Server{Store: store}
	r := gin.Default()
	r.GET("/get-house/:id", server.GetHouseByID)
	r.GET("/get-all-houses", server.GetAllHouses)
	r.POST("/add-new-house", server.AddNewHouse)
	r.POST("/add-new-image/:id", server.UpdateNewImage)
	r.GET("/upload-image/:id", server.GetImageByID)
	if err := r.Run("localhost:8080"); err != nil {
		return
	}
}
