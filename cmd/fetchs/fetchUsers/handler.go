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
	GetAllUsers(rootUser string) (map[string]map[string]string, map[string][]string, error)
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	users := make(map[string]map[string]string)
	queryNametoUserId := make(map[string][]string)
	var err error
	rootUser, exists := request.PathParameters["rootUser"]
	if !exists {
		log.Println("NO ROOT USER  USER PROVIDED WHEN QUERYING USERS")
		return responses.BadRequest(), nil
	}
	user, exists := request.PathParameters["user"]
	if !exists {
		log.Println("NO USER PROVIDED WHEN QUERYING USERS")
		return responses.BadRequest(), nil
	}

	if user == "all" {
		users, queryNametoUserId, err = h.dbClient.GetAllUsers(rootUser)
		if err != nil {
			log.Println("ERROR WHEN CALLING db.GetAllUsers")
			return responses.ServerError(err), fmt.Errorf("ERROR WHEN CALLING db.GetAllUsers")
		}
	} else {
		users = make(map[string]map[string]string)
		queryNametoUserId = make(map[string][]string)
	}
	rawResponseBody := map[string]interface{}{
		"users":             users,
		"queryNameToUserId": queryNametoUserId,
	}
	responseBody, err := json.Marshal(rawResponseBody)
	if err != nil {
		log.Println("ERROR MARSHALLING RESPONSE BODY TO JSON")
		return responses.ServerError(err), fmt.Errorf("ERROR MARSHALLING RESPONSE BODY TO JSON")
	}

	// request.Headers["Access-Control-Allow-Origin"] = "*"
	// request.Headers["Access-Control-Allow-Headers"] = "*"
	// request.Headers["Access-Control-Allow-Credentials"] = "true"

	log.Println("FETCH USERS | fetch artists lambda successful exectuion")
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
