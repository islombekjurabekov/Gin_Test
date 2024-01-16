package rest

import (
	"Gin_test/database"
	"Gin_test/domain"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Server struct {
	Store database.DBStore
}

func (s *Server) GetHouseByID(c *gin.Context) {
	id := c.Param("id")
	res, err := s.Store.GetHouseByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) GetAllHouses(c *gin.Context) {
	res, err := s.Store.GetAllHouses()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) AddNewHouse(c *gin.Context) {
	var req domain.HomeRest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	s.Store.AddNewHouse(req)
	c.JSON(http.StatusOK, gin.H{
		"Message": "New Information successfully added to database",
	})
}

func (s *Server) UpdateNewImage(c *gin.Context) {
	id := c.Param("id")
	var buf = bytes.Buffer{}

	image, res, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	if _, err = io.Copy(&buf, image); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	s.Store.UpdateNewImage(buf, id, res)
	c.JSON(http.StatusOK, gin.H{
		"Message": "Image successfully added",
	})
}

func (s *Server) GetImageByID(c *gin.Context) {
	id := c.Param("id")

	fileName, err := s.Store.GetImageByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Set("Content-Type", "multipart/form-data")
	c.File("./picture/" + fileName)

}
