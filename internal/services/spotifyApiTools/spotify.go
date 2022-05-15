package spotifyApiTools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/MowMowchow/musicial/internal/constants"
	"github.com/MowMowchow/musicial/internal/helpers"
	"github.com/MowMowchow/musicial/internal/models"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c Client) GetAllUserPlaylists(userId string, accessToken string) (map[string]models.SpotifyPlaylist, error) {
	var playlistArr []models.SpotifyPlaylist

	req, err := http.NewRequest("GET", constants.SpotifyUserPlaylistsURL, nil)
	if err != nil {
		log.Println("error building base request for fetching user playlists")
		return nil, fmt.Errorf("error building base request for fetching user playlists")
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	offset := 0
	limit := 20
	total := 0
	urlQueries := req.URL.Query()
	urlQueries.Add("offset", helpers.ToString(offset))
	urlQueries.Add("limit", helpers.ToString(limit))

	// initial playlist fetch
	resp, err := c.HttpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("error playlist retrieval ", err)
		return nil, fmt.Errorf("error playlist retrieval ")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error with ioutil.ReadAll(resp.Body)", err)
		return nil, fmt.Errorf("Error with ioutil.ReadAll(resp.Body)")
	}

	var userPlaylistsResponseBody models.SpotifyUserPlaylistsResponseBody
	err = json.Unmarshal(respBody, &userPlaylistsResponseBody)
	if err != nil {
		log.Println("error unmarshalling response body from spotify get user playlists", err)
		return nil, fmt.Errorf("error unmarshalling response body from spotify get user playlists")
	}

	for _, playlist := range userPlaylistsResponseBody.Items {
		playlistArr = append(playlistArr, playlist)
	}

	offset += limit
	total = userPlaylistsResponseBody.Total

	// if the user has >= 20 playlists
	var errorWithFetch error = nil
	var wg sync.WaitGroup

	for ; offset < total; offset += limit {
		wg.Add(1)
		go func(pageOffset int) {
			defer wg.Done()
			req, err := http.NewRequest("GET", constants.SpotifyUserPlaylistsURL, nil)
			if err != nil {
				log.Println("error building base request for fetching user playlists")
				errorWithFetch = err
			}

			req.Header.Add("Authorization", "Bearer "+accessToken)
			req.Header.Add("Content-Type", "application/json")

			urlQueries := req.URL.Query()
			urlQueries.Add("offset", helpers.ToString(pageOffset))
			urlQueries.Add("limit", helpers.ToString(limit))
			req.URL.RawQuery = urlQueries.Encode()

			resp, err := c.HttpClient.Do(req)
			if err != nil || resp.StatusCode != 200 {
				log.Println("error playlist retrieval | offset: ", pageOffset)
				errorWithFetch = err
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error with ioutil.ReadAll(resp.Body)")
			}

			var userPlaylistsResponseBody models.SpotifyUserPlaylistsResponseBody
			err = json.Unmarshal(respBody, &userPlaylistsResponseBody)
			if err != nil {
				log.Println("error unmarshalling response body from spotify get user playlists")
				errorWithFetch = err
			}

			if len(userPlaylistsResponseBody.Items) > 0 {
				for _, playlist := range userPlaylistsResponseBody.Items {
					playlistArr = append(playlistArr, playlist)
				}
			}
		}(offset)
	}

	wg.Wait()

	if (errorWithFetch != nil && len(playlistArr) != total) || (len(playlistArr) != total) {
		return nil, errorWithFetch
	}

	playlists := make(map[string]models.SpotifyPlaylist)
	for _, playlist := range playlistArr {
		playlists[playlist.Id] = playlist
	}

	return playlists, nil
}

// func (c Client) GetAllUserLikedSongs(userId string, accessToken string) (map[string]string, error) {

// }

func (c Client) GetAllPlaylistTracks(userId string, accessToken string, playlistId string, totalTracks int) (map[string]models.SpotifyTrack, error) {

	url := fmt.Sprintf(constants.SpotifyPlaylistTracks, playlistId)
	limit := 100
	var errorWithFetch error = nil

	var wg sync.WaitGroup

	var tracksArr []models.SpotifyTrack
	tracks := make(map[string]models.SpotifyTrack)
	for offset := 0; offset < totalTracks; offset += limit {
		wg.Add(1)
		go func(pageOffset int) {
			defer wg.Done()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println("error building base request for fetching user playlists")
				errorWithFetch = err
			}

			req.Header.Add("Authorization", "Bearer "+accessToken)
			req.Header.Add("Content-Type", "application/json")

			urlQueries := req.URL.Query()
			urlQueries.Add("offset", helpers.ToString(pageOffset))
			urlQueries.Add("limit", helpers.ToString(limit))
			urlQueries.Add("fields", constants.SpotifyGetPlaylistItemsFields)
			req.URL.RawQuery = urlQueries.Encode()

			resp, err := c.HttpClient.Do(req)
			if err != nil || resp.StatusCode != 200 {
				log.Println("error playlist tracks retrieval | offset: "+helpers.ToString(pageOffset)+"| err: ", err)
				errorWithFetch = err
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error with ioutil.ReadAll(resp.Body)")
			}

			var userPlaylistTracksResponseBody models.SpotifyPlaylistTracksResponseBody
			err = json.Unmarshal(respBody, &userPlaylistTracksResponseBody)
			if err != nil {
				errorWithFetch = err
			}

			if len(userPlaylistTracksResponseBody.Items) > 0 {
				for _, track := range userPlaylistTracksResponseBody.Items {
					tracksArr = append(tracksArr, track.Track)
				}
			}
		}(offset)
	}

	wg.Wait()

	if (errorWithFetch != nil && len(tracksArr) != totalTracks) || (len(tracksArr) != totalTracks) {
		return nil, errorWithFetch
	}

	for _, track := range tracksArr {
		tracks[track.Id] = track
	}
	return tracks, nil

}

func (c Client) RefreshToken(refreshToken string, clientId string, clientSecret string) (string, string, error) {
	req, err := http.NewRequest("POST", constants.SpotifyTokenURL, nil)
	if err != nil {
		log.Println("error building base request for callback refresh access token retrieval | ", err)
		return "", "", fmt.Errorf("error building base request for callback refresh access token retrieval")
	}

	queries := req.URL.Query()
	queries.Add("grant_type", "refresh_token")
	queries.Add("refresh_token", refreshToken)
	queries.Add("json", "true")
	req.URL.RawQuery = queries.Encode()

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString((bytes.NewBufferString(clientId+":"+clientSecret)).Bytes()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HttpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("error refresh access token retrieval | ", err)
		return "", "", fmt.Errorf("error refresh access token retrieval ")
	}

	var tokenResponseBody models.SpotifyTokenResponseBody
	err = json.NewDecoder(resp.Body).Decode(&tokenResponseBody)
	if err != nil {
		log.Println("error decoding response body from spotify refresh access token post request | ", err)
		return "", "", fmt.Errorf("error decoding response body from spotify refresh access token post request")
	}

	return tokenResponseBody.AccessToken, tokenResponseBody.RefreshToken, nil
}
