package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MowMowchow/musicial/internal/responses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	httpClient  HttpClient
	redisClient RedisClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RedisClient interface {
	HGet(ctx context.Context, key, field string) *redis.StringCmd
}

var ctx = context.Background()

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("PROCESS LOGIN | starting get user process status")
	userId, exists := request.PathParameters["userId"]
	if !exists {
		log.Println("pollLogin | no userId was provided in path parameters")
		return responses.ServerError(), fmt.Errorf("pollLogin | no userId was provided in path parameters")

	}
	isProcessingVal := 1

	isProcessing, err := h.redisClient.HGet(ctx, userId, "isProcessing").Result()
	if err != nil {
		return responses.ServerError(err), fmt.Errorf("error setting user process status")
	}
	if err != nil {
		log.Println("error setting user proccess status | ", err)
	}
	log.Println("PROCESS LOGIN | get user process status successful")

	if isProcessing == "0" {
		isProcessingVal = 0
	} else if isProcessing != "1" {
		log.Println("pollLogin | error with process login | isProcessing: ", isProcessing)
		isProcessingVal = 2
	}

	rawResponseBody := map[string]int{
		"isProcessing": isProcessingVal,
	}

	responseBody, err := json.Marshal(rawResponseBody)
	if err != nil {
		log.Println("ERROR MARSHALLING RESPONSE BODY TO JSON", err)
		return responses.ServerError(err), fmt.Errorf("ERROR MARSHALLING RESPONSE BODY TO JSON")
	}

	log.Println("PROCESS LOGIN | get user process status successful | rawResponseBody: ", rawResponseBody)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(responseBody),
	}, nil
}
