package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	sLambda "github.com/aws/aws-sdk-go/service/lambda"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
)

func main() {
	httpClient := &http.Client{}

	redisClientAddress := os.Getenv("REDIS_CLIENT_ADDRESS")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisClientAddress,
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	lambdaClient := sLambda.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	handler := Handler{
		httpClient:   httpClient,
		redisClient:  redisClient,
		lambdaClient: lambdaClient,
	}

	lambda.Start(handler.HandleRequest)

}
