package auth

import (
	"context"
	"log"
	"time"
)

// return Bool(whether or not the user is valid), string(newAccessToken), error
func (c Client) IsValid(ctx context.Context, userId string, accessToken string, clientId string, clientSecret string) (bool, string, error) {
	userInfo, err := c.RedisClient.HGetAll(ctx, userId).Result()
	if err != nil {
		return false, "", err
	}
	if len(userInfo) > 0 {
		accessTokenExpires, err := time.Parse(time.RFC3339, userInfo["accessTokenExpires"])
		if err != nil {
			return false, "", err
		}

		if accessTokenExpires.Add(-25 * time.Minute).After(time.Now().UTC()) { // if there are >= 30min left to session
			if storedAccessToken, exists := userInfo["accessToken"]; exists {
				if storedAccessToken == accessToken {
					return true, "", nil
				}
			}
		}
		if storedRefreshToken, exists := userInfo["refreshToken"]; exists {
			newAccessToken, newRefreshToken, err := c.SpotifyClient.RefreshToken(storedRefreshToken, clientId, clientSecret)
			if err != nil {
				return false, "", err
			}
			newAccessTokenExpires := time.Now().UTC().Add(55 * time.Minute).Format(time.RFC3339)
			redisHash := make(map[string]interface{})
			if newRefreshToken == "" {
				redisHash["accessToken"] = newAccessToken
				redisHash["accessTokenExpires"] = newAccessTokenExpires
			} else {
				redisHash["accessToken"] = newAccessToken
				redisHash["refreshToken"] = newRefreshToken
				redisHash["accessTokenExpires"] = newAccessTokenExpires
			}
			err = c.RedisClient.HSet(
				ctx,
				userId,
				redisHash,
			).Err()
			if err != nil {
				log.Println("error caching user info | ", err)
				return false, "", err
			}
			return true, newAccessToken, nil
		}
	}
	return false, "", nil
}
