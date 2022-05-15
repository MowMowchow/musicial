package db

import (
	"net/http"

	"github.com/MowMowchow/musicial/internal/services/spotifyApiTools"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Client struct {
	Neo4jClient neo4j.Session
	HttpClient  http.Client
	// SpotifyApiClient spotifyApiTools.Client{HttpClient: *httpClient,}
	SpotifyClient spotifyApiTools.Client
}
