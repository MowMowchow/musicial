package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MowMowchow/musicial/internal/services/auth"
	"github.com/MowMowchow/musicial/internal/services/db"
	"github.com/MowMowchow/musicial/internal/services/spotifyApiTools"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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

	// neo4j setup
	neo4jClientAddress := os.Getenv("NEO4J_CLIENT_ADDRESS")
	neo4jClientUser := os.Getenv("NEO4J_CLIENT_USER")
	neo4jClientPassword := os.Getenv("NEO4J_CLIENT_PASSWORD")
	neo4jDriver, err := neo4j.NewDriver(neo4jClientAddress, neo4j.BasicAuth(neo4jClientUser, neo4jClientPassword, ""))
	if err != nil {
		log.Println("FAILED TO CREATE NEO4J DRIVER")
	}
	defer neo4jDriver.Close()
	neo4jClient := neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer neo4jClient.Close()

	dbClient := db.Client{
		HttpClient:    *httpClient,
		Neo4jClient:   neo4jClient,
		SpotifyClient: spotifyApiTools.Client{HttpClient: *httpClient},
	}
	authClient := auth.Client{
		RedisClient:   redisClient,
		SpotifyClient: spotifyApiTools.Client{HttpClient: *httpClient},
	}

	handler := Handler{
		httpClient: httpClient,
		dbClient:   dbClient,
		authClient: authClient,
	}

	lambda.Start(handler.HandleRequest)

}
