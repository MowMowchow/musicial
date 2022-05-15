package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MowMowchow/musicial/internal/models"
	"github.com/MowMowchow/musicial/internal/responses"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	httpClient HttpClient
	dbClient   DbClient
	authClient AuthClient
}

type HttpClient interface {
}

type DbClient interface {
	GetCommonArtists(userId string) (map[string]*models.MusicialSearchArtistResult, map[string][]string, error)
}

type AuthClient interface {
	IsValid(ctx context.Context, userId string, accessToken string, clientId string, clientSecret string) (bool, string, error)
}

var ctx = context.Background()

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	var err error
	userId, exists := request.PathParameters["user"]
	if !exists {
		log.Println("NO userId PROVIDED WHEN QUERYING USERS")
		return responses.BadRequest(), nil
	}

	var requestBody map[string]string
	err = json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Println("error unmarshalling response body from udpate user request | ", err)
		return responses.ServerError(err), fmt.Errorf("error unmarshalling response body from spotify search artist request")
	}

	accessToken, exists := requestBody["accessToken"]
	if !exists {
		log.Println("NO ACCESS TOKEN PROVIDED IN HEADERS WHEN ATTEMPTING TO CONNECT")
		return responses.BadRequest(), nil
	}

	// call auth shit

	isValid, newAccessToken, err := h.authClient.IsValid(ctx, userId, accessToken, clientID, clientSecret)
	log.Println("AUTH IS VALID: ", isValid, " | NEW ACCESS TOKEN? : ", newAccessToken)

	artists, queryNametoArtistId, err := h.dbClient.GetCommonArtists(userId)
	if err != nil {
		log.Println("ERROR WHEN CALLING db.getCommonArtists", err)
		return responses.ServerError(err), fmt.Errorf("ERROR WHEN CALLING db.getCommonArtists")
	}

	rawResponseBody := map[string]interface{}{
		"artists":             artists,
		"queryNameToArtistId": queryNametoArtistId,
	}
	responseBody, err := json.Marshal(rawResponseBody)
	if err != nil {
		log.Println("ERROR MARSHALLING RESPONSE BODY TO JSON", err)
		return responses.ServerError(err), fmt.Errorf("ERROR MARSHALLING RESPONSE BODY TO JSON")
	}

	log.Println("FETCH ARTISTS | fetch artists lambda successful exectuion")
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
