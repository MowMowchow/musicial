package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/MowMowchow/musicial/internal/constants"
	"github.com/MowMowchow/musicial/internal/models"
	"github.com/MowMowchow/musicial/internal/responses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	httpClient  HttpClient
	redisClient RedisClient
	dbClient    DbClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RedisClient interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

type DbClient interface {
	GetPlaylistsToUpdate(userId string, accessToken string) ([]models.SpotifyPlaylist, error)
	UpdatePlaylists(userId string, accessToken string) error
	UpdateUser(user models.SpotifyUserProfileResponseBody, accessToken string) error
}

var ctx = context.Background()

func (h *Handler) updateProcessStateToError(userId string) error {
	err := h.redisClient.HSet(
		ctx,
		userId,
		map[string]interface{}{
			"isProcessing": "2",
		},
	).Err()
	if err != nil {
		log.Println("error setting user proccess status to error| ", err)
		// return responses.ServerError(err), fmt.Errorf("error setting user process status")
		return fmt.Errorf("error setting user process status to error")
	}
	return nil
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var tokenResponseBody models.SpotifyTokenResponseBodyWithUserId
	err := json.Unmarshal([]byte(request.Body), &tokenResponseBody)
	if err != nil {
		log.Println("error unmarshalling response body from process login base request | ", err)
		// return fmt.Errorf("error unmarshalling response body from process login base request")
		return responses.ServerError(err), fmt.Errorf("error unmarshalling response body from process login base request")
	}

	log.Println("PROCESS LOGIN | Sending request for user profile retrieval")
	req, err := http.NewRequest("GET", constants.SpotifyUserURL, nil)
	req.Header.Add("Authorization", "Bearer "+tokenResponseBody.AccessToken)
	resp, err := h.httpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("error with get request response for user profile retrieval | ", err, " | statusCode: ", resp.StatusCode)
		h.updateProcessStateToError(tokenResponseBody.UserId)
		return responses.ServerError(err), fmt.Errorf("error with get request response for user profile retrieval")
		// return fmt.Errorf("error with get request response for user profile retrieval")
	}
	log.Println("PROCESS LOGIN | Request for user profile retrieval successful | resp: ", resp)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error with ioutil.ReadAll(resp.Body) | ", err)
		h.updateProcessStateToError(tokenResponseBody.UserId)
		return responses.ServerError(err), fmt.Errorf("Error with ioutil.ReadAll(resp.Body)")
		// return fmt.Errorf("Error with ioutil.ReadAll(resp.Body)")
	}
	var userProfileResponseBody models.SpotifyUserProfileResponseBody
	err = json.Unmarshal(respBody, &userProfileResponseBody)
	if err != nil {
		log.Println("error unmarshalling response body from spotify user profile get request | ", err)
		h.updateProcessStateToError(tokenResponseBody.UserId)
		return responses.ServerError(err), fmt.Errorf("error unmarshalling response body from spotify user profile get request")
		// return fmt.Errorf("error unmarshalling response body from spotify user profile get request")
	}

	processedUserInformation := true
	err = h.dbClient.UpdateUser(userProfileResponseBody, tokenResponseBody.AccessToken)
	if err != nil {
		log.Println("ERROR UPDATING USER", err)
		processedUserInformation = false
	}
	log.Println("PROCESS LOGIN | update user successful")

	err = h.dbClient.UpdatePlaylists(userProfileResponseBody.Id, tokenResponseBody.AccessToken)
	if err != nil {
		log.Println("ERROR UPDATING PLAYLISTS", err)
		processedUserInformation = false
	}
	log.Println("PROCESS LOGIN | update playlists successful")

	respBody, err = json.Marshal(map[string]bool{"processedUserInformation": processedUserInformation})
	if err != nil {
		h.updateProcessStateToError(tokenResponseBody.UserId)
		return responses.ServerError(err), fmt.Errorf("error marshalling processedUserInformation boolean for response body")
		// return fmt.Errorf("error marshalling processedUserInformation boolean for response body")
	}

	err = h.redisClient.HSet(
		ctx,
		userProfileResponseBody.Id,
		map[string]interface{}{
			"isProcessing": "0",
		},
	).Err()
	if err != nil {
		log.Println("error setting user proccess status | ", err)
		h.updateProcessStateToError(tokenResponseBody.UserId)
		return responses.ServerError(err), fmt.Errorf("error setting user process status")
		// return fmt.Errorf("error setting user process status")
	}
	log.Println("PROCESS LOGIN | set user process status successful")

	log.Println("PROCESS LOGIN | PROCESS LOGIN SUCCESS")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
	// return nil
}
