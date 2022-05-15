package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MowMowchow/musicial/internal/constants"
	"github.com/MowMowchow/musicial/internal/helpers"
	"github.com/MowMowchow/musicial/internal/responses"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	responseType := "code"
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	scope := "user-follow-read user-read-recently-played playlist-read-collaborative user-library-read user-library-read user-follow-read user-library-modify playlist-modify-public user-top-read "
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	state := helpers.CreateRandomString(16)

	req, err := http.NewRequest("GET", constants.BaseSpotifyAuthUrl, nil)
	if err != nil {
		return responses.ServerError(), fmt.Errorf("error building the base request")
	}

	queryURL := req.URL.Query()
	queryURL.Add("response_type", responseType)
	queryURL.Add("client_id", clientID)
	queryURL.Add("scope", scope)
	queryURL.Add("redirect_uri", redirectURI)
	queryURL.Add("state", state)
	encodedQueryURL := queryURL.Encode()

	request.Headers["location"] = constants.BaseSpotifyAuthUrl + encodedQueryURL
	request.Headers["Set-Cookie"] = constants.StateKey + "=" + state

	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"location":                         constants.BaseSpotifyAuthUrl + encodedQueryURL,
			"Set-Cookie":                       constants.StateKey + "=" + state,
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: request.Body,
	}, nil
}
