package ConnectDatabase

import (
	"Gin_test/database"
	"Gin_test/rest"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func Connection() (database.DBStore, rest.RedisC) {
	//Connect to Database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	//Connect to Redis

	var pong string
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to Redis: ", pong)
	}
	rClient := rest.Client(client)
	store := database.New(DB)

	return store, rClient
}
