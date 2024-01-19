package rest

import (
	"Gin_test/database"
	"Gin_test/domain"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"time"
)

type RedisC struct {
	RClient *redis.Client
}

func Client(cRedis *redis.Client) RedisC {
	return RedisC{
		RClient: cRedis,
	}
}

type Server struct {
	Store database.DBStore
	Redis RedisC
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

func (s *Server) SaveToRedisFile(c *gin.Context) {
	var buf bytes.Buffer
	image, res, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if _, err = io.Copy(&buf, image); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	_, err = s.Redis.RClient.Set(context.Background(), res.Filename, buf.String(), 2*time.Minute).Result()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"key": res.Filename,
	})
}

func (s *Server) AddNewHouse(c *gin.Context) {
	var (
		req   domain.HomeRest
		image string
	)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	image, err = s.Redis.RClient.Get(context.Background(), req.ImageKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	req.Image = image
	req.ImageName = req.ImageKey
	s.Store.AddNewHouse(req)
	c.JSON(http.StatusOK, gin.H{
		"Message": "Data successfully added to database",
	})
}

func (s *Server) GetImageByID(c *gin.Context) {
	id := c.Param("id")
	fileName, err := s.Store.GetImageByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.Writer.Header().Set("Content-Disposition", "attachment; filename = "+fileName)
	c.Writer.Header().Set("Content-Type", "multipart/form-data")
	c.File("./picture/" + fileName)

}
