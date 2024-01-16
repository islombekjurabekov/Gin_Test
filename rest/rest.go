package rest

import (
	"Gin_test/database"
	"Gin_test/domain"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func GetHouseByID(c *gin.Context) {
	id := c.Param("id")
	res, err := database.GetHouseByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, res)
}

func GetAllHouses(c *gin.Context) {
	res, err := database.GetAllHouses()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, res)
}

func AddNewHouse(c *gin.Context) {
	var req domain.HomeRest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	database.AddNewHouse(req)
	c.JSON(http.StatusOK, gin.H{
		"Message": "New Information successfully added to database",
	})
}

func UpdateNewImage(c *gin.Context) {
	id := c.Param("id")
	var buf = bytes.Buffer{}

	image, res, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	if _, err = io.Copy(&buf, image); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	database.UpdateNewImage(buf, id, res)
	c.JSON(http.StatusOK, gin.H{
		"Message": "Image successfully added",
	})
}

func GetImageByID(c *gin.Context) {
	id := c.Param("id")

	err := database.GetImageByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+"1.png")
	c.Writer.Header().Set("Content-Type", "multipart/form-data")
	c.File("./picture/Big_house.png")
}
