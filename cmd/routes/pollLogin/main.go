package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
)

func main() {
	httpClient := &http.Client{}

	// redis client setup
	redisClientAddress := os.Getenv("REDIS_CLIENT_ADDRESS")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisClientAddress,
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	handler := Handler{
		httpClient:  httpClient,
		redisClient: redisClient,
	}

	lambda.Start(handler.HandleRequest)

}
