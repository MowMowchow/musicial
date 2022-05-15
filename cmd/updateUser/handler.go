package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MowMowchow/musicial/internal/responses"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	httpClient HttpClient
	dbClient   DbClient
}

type HttpClient interface {
}

type DbClient interface {
	ConnectUsers(rootUserId string, connectUserID string) (bool, error)
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody map[string]string
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Println("error unmarshalling response body from udpate user request | ", err)
		return responses.ServerError(err), fmt.Errorf("error unmarshalling response body from spotify update user request")
	}

	rootUserId, exists := requestBody["userId"]
	if !exists {
		log.Println("NO USER ID PROVIDED IN HEADERS WHEN ATTEMPTING TO CONNECT")
		return responses.BadRequest(), nil
	}

	// accessToken, exists := responseBody["accessToken"]
	// if !exists {
	// 	log.Println("NO ACCESS TOKEN PROVIDED IN HEADERS WHEN ATTEMPTING TO CONNECT")
	// 	return responses.BadRequest(), nil
	// }

	// call auth shit

	connectUserId, exists := request.PathParameters["user"]
	if !exists {
		log.Println("NO CONNECT USER ID PROVIDED IN PATH PARAMETER WHEN ATTEMPTING TO CONNECT")
		return responses.BadRequest(), nil
	}

	dbResult, err := h.dbClient.ConnectUsers(rootUserId, connectUserId)
	if err != nil {
		log.Println("ERROR WHEN CALLING db.ConnectUsers")
		return responses.ServerError(err), fmt.Errorf("ERROR WHEN CALLING db.ConnectUsers")
	}

	responseBody, err := json.Marshal(dbResult)
	if err != nil {
		log.Println("ERROR MARSHALLING RESPONSE BODY TO JSON")
		return responses.ServerError(err), fmt.Errorf("ERROR MARSHALLING RESPONSE BODY TO JSON")
	}

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
