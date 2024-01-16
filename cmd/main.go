package main

import (
	"Gin_test/ConnectDatabase"
	"Gin_test/rest"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ConnectDatabase.Connection()
	r := gin.Default()
	r.GET("/get-house/:id", rest.GetHouseByID)
	r.GET("/get-all-houses", rest.GetAllHouses)
	r.POST("/add-new-house", rest.AddNewHouse)
	r.POST("/add-new-image/:id", rest.UpdateNewImage)
	r.GET("/upload-image/:id", rest.GetImageByID)
	if err := r.Run("localhost:8080"); err != nil {
		return
	}
}
