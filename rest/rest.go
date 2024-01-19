package rest

import (
	"Gin_test/database"
	"Gin_test/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s *Server) AddNewHouseByKeyToDatabase(c *gin.Context) {
	key := c.Param("key")
	ConvertedInfo, err := s.Redis.RClient.Get(context.Background(), key).Bytes()
	if err != nil {
		log.Fatal(err)
	}
	var param domain.HomeRest
	err = json.Unmarshal(ConvertedInfo, &param)
	s.Store.AddNewHouseByKeyToDatabase(param)
	c.JSON(http.StatusOK, gin.H{
		"Message": "New Information successfully added to database",
	})
}

func (s *Server) AddNewHouse(c *gin.Context) {
	var req domain.HomeRest
	var res string
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	key := uuid.New().String()
	ConvertedInfo, _ := json.Marshal(req)
	res, err = s.Redis.RClient.Set(context.Background(), key, ConvertedInfo, 2*time.Minute).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully saved to Redis storage:", res)
	c.JSON(http.StatusOK, gin.H{
		"key": key,
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
	c.Writer.Header().Set("Content-Disposition", "attachment; filename = "+fileName)
	c.Writer.Header().Set("Content-Type", "multipart/form-data")
	c.File("./picture/" + fileName)

}
