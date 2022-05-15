package auth

import (
	"github.com/MowMowchow/musicial/internal/services/spotifyApiTools"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	RedisClient   *redis.Client
	SpotifyClient spotifyApiTools.Client
}
