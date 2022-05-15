package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MowMowchow/musicial/internal/constants"
	"github.com/MowMowchow/musicial/internal/models"
	"github.com/MowMowchow/musicial/internal/responses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redis/redis/v8"

	"github.com/aws/aws-sdk-go/aws"
	sLambda "github.com/aws/aws-sdk-go/service/lambda"
)

type Handler struct {
	httpClient   HttpClient
	redisClient  RedisClient
	lambdaClient LambdaClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RedisClient interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

type LambdaClient interface {
	Invoke(input *sLambda.InvokeInput) (*sLambda.InvokeOutput, error)
}

var ctx = context.Background()

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	code := request.QueryStringParameters["code"]
	state := request.QueryStringParameters["state"]

	log.Println("PROCESS LOGIN | recieved request: ", request)

	var storedState string
	if storedStateCookie, valid := request.Headers["Cookie"]; valid {
		log.Println("HEADERS COOKIE: ", request.Headers["Cookie"])
		storedState = storedStateCookie
		delete(request.Headers, "Cookie")
	}

	if state == "" || constants.StateKey+"="+state != storedState {
		log.Println("state_mismatch error")
		return responses.BadRequest(), fmt.Errorf("state_mismatch error")
	}

	// getting tokens
	req, err := http.NewRequest("POST", constants.SpotifyTokenURL, nil)
	if err != nil {
		log.Println("error building base request for callback token retrieval | ", err)
		return responses.ServerError(err), fmt.Errorf("error building base request for callback token retrieval")
	}

	queries := req.URL.Query()
	queries.Add("code", code)
	queries.Add("redirect_uri", redirectURI)
	queries.Add("grant_type", "authorization_code")
	queries.Add("json", "true")
	req.URL.RawQuery = queries.Encode()

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString((bytes.NewBufferString(clientID+":"+clientSecret)).Bytes()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Println("PROCESS LOGIN | Sending request for user token retrieval")
	tokenResp, err := h.httpClient.Do(req)
	if err != nil || tokenResp.StatusCode != 200 {
		log.Println("error token retrieval | ", err)
		return responses.ServerError(err), fmt.Errorf("error token retrieval ")
	}
	log.Println("PROCESS LOGIN | Request for user token retrieval successful")

	var tokenResponseBody models.SpotifyTokenResponseBodyWithUserId
	err = json.NewDecoder(tokenResp.Body).Decode(&tokenResponseBody)
	if err != nil {
		log.Println("error decoding response body from spotify token post request | ", err)
		return responses.ServerError(err), fmt.Errorf("error decoding response body from spotify token post request")
	}

	// setting up request for user profile retrieval

	log.Println("PROCESS LOGIN | Sending request for user profile retrieval")
	req, err = http.NewRequest("GET", constants.SpotifyUserURL, nil)
	req.Header.Add("Authorization", "Bearer "+tokenResponseBody.AccessToken)
	resp, err := h.httpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("error with get request response for user profile retrieval | ", err, " | statusCode: ", resp.StatusCode)
		return responses.ServerError(err), fmt.Errorf("error with get request response for user profile retrieval")
	}
	log.Println("PROCESS LOGIN | Request for user profile retrieval successful | resp: ", resp)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error with ioutil.ReadAll(resp.Body) | ", err)
		return responses.ServerError(err), fmt.Errorf("Error with ioutil.ReadAll(resp.Body)")
	}
	var userProfileResponseBody models.SpotifyUserProfileResponseBody
	err = json.Unmarshal(respBody, &userProfileResponseBody)
	if err != nil {
		log.Println("error unmarshalling response body from spotify user profile get request | ", err)
		return responses.ServerError(err), fmt.Errorf("error unmarshalling response body from spotify user profile get request")
	}

	// set redis info (start processing)

	if tokenResponseBody.AccessToken != "" && tokenResponseBody.RefreshToken != "" && tokenResponseBody.ExpiresIn != 0 {
		expiresTime := time.Now().UTC().Add(55 * time.Minute).Format(time.RFC3339)
		err = h.redisClient.HSet(
			ctx,
			userProfileResponseBody.Id,
			map[string]interface{}{
				"accessToken":        tokenResponseBody.AccessToken,
				"refreshToken":       tokenResponseBody.RefreshToken,
				"accessTokenExpires": expiresTime,
				"isProcessing":       "1",
			},
		).Err()
		if err != nil {
			log.Println("error caching user info | ", err)
			return responses.ServerError(err), fmt.Errorf("error caching user info")
		}
	} else {
		log.Println("token response did not contain required information")
		return responses.ServerError(), fmt.Errorf("token response did not contain required information")
	}
	log.Println("HANDLE REDIRECT | redis user successful")

	// setting up process login request
	log.Println("TOKEN RESPONSE BODY WHEN INVOKING PROCESS LOGIN: ", tokenResponseBody)
	tokenResponseBody.UserId = userProfileResponseBody.Id
	responseBody, err := json.Marshal(tokenResponseBody)
	if err != nil {
		log.Println("ERROR MARSHALLING RESPONSE BODY TO JSON", err)
		return responses.ServerError(err), fmt.Errorf("ERROR MARSHALLING RESPONSE BODY TO JSON")
	}

	lambdaPayload := models.LambdaPayload{
		Body: string(responseBody),
	}
	jsonLambdaPayload, err := json.Marshal(lambdaPayload)
	if err != nil {
		log.Println("ERROR MARSHALLING LAMBDA PAYLOAD TO JSON", err)
		return responses.ServerError(err), fmt.Errorf("ERROR MARSHALLING LAMBDA PAYLOAD TO JSON")
	}

	result, err := h.lambdaClient.Invoke(&sLambda.InvokeInput{
		FunctionName:   aws.String("musicial-dev-processLogin"),
		InvocationType: aws.String("Event"),
		Payload:        jsonLambdaPayload,
	})
	if err != nil {
		log.Println("ERROR INVOKING PROCESS LOGIN", err)
		return responses.ServerError(err), fmt.Errorf("ERROR INVOKING PROCESS LOGIN")
	}

	log.Println("HANDLE REDIRECT | successfully invoked request for process login | result ", result)

	log.Println("PROCESS LOGIN | handle redirect was successful")
	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"location": "http://musicial.net/pollLogin/?userId=" + userProfileResponseBody.Id + "&accessToken=" + tokenResponseBody.AccessToken,
		},
		MultiValueHeaders: map[string][]string{
			"Set-Cookie": {"accessToken " + tokenResponseBody.AccessToken, "userId " + tokenResponseBody.UserId},
		},
	}, nil
}
