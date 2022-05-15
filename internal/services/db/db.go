package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/MowMowchow/musicial/internal/helpers"
	"github.com/MowMowchow/musicial/internal/models"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func (c Client) Neo4jWriteTransaction(query string, queryParams map[string]interface{}) (interface{}, error) {
	result, err := c.Neo4jClient.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		runResult, err := transaction.Run(query, queryParams)
		if err != nil {
			return nil, err
		}
		return nil, runResult.Err()
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Returns a slice of 'row' slices
func (c Client) Neo4jReadTransaction(query string, queryParams map[string]interface{}) ([][]interface{}, error) {
	var rows [][]interface{}
	_, err := c.Neo4jClient.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		runResult, err := transaction.Run(query, queryParams)
		if err != nil {
			return nil, err
		}
		if runResult != nil {
			dbRecords, err := runResult.Collect()
			if err != nil {
				return nil, err
			}
			for _, row := range dbRecords {
				rows = append(rows, row.Values)
			}
			return rows, runResult.Err()
		}

		return nil, runResult.Err()
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c Client) MapPlaylistIdToSnapshotIdFromResult(results [][]interface{}) (map[string]map[string]string, error) {
	if results == nil {
		log.Println("results is null")
		return nil, nil
	}
	playlistIdToSnapshotId := make(map[string]map[string]string)
	for _, row := range results {
		if row[0] == nil {
			row[0] = ""
		}
		if row[1] == nil {
			row[1] = ""
		}
		if row[2] == nil {
			row[2] = ""
		}
		playlistId, isString1 := row[0].(string)
		totalTracks, isString2 := row[1].(string)
		snapshotId, isString3 := row[2].(string)
		if isString1 && isString2 && isString3 {
			playlistIdToSnapshotId[playlistId] = make(map[string]string)
			playlistIdToSnapshotId[playlistId]["snapshotId"] = snapshotId
			playlistIdToSnapshotId[playlistId]["totalTracks"] = totalTracks
			playlistIdToSnapshotId[playlistId]["exists"] = "0"
		} else {
			return nil, fmt.Errorf("Error when getting playlist snapshot_id from result: row value is not string ")
		}
	}
	return playlistIdToSnapshotId, nil
}

func (c Client) MapTrackIdToTrackNameFromResult(results [][]interface{}) (map[string]string, error) {
	trackIdToTrackName := make(map[string]string)
	for _, row := range results {
		trackId, isString1 := row[0].(string)
		trackName, isString2 := row[1].(string)
		if isString1 && isString2 {
			trackIdToTrackName[trackId] = trackName
		} else {
			log.Println("MapTrackIdTrackNameFromResult SOMTHING IN THE ROW WAS NOT A STRING: ", row)
			return nil, fmt.Errorf("Error when getting songId from result")
		}
	}
	return trackIdToTrackName, nil
}

func (c Client) GetBatchInsertRows(userId string, tracks map[string]models.SpotifyTrack, playlist models.SpotifyPlaylist) []map[string]string {
	rows := make([]map[string]string, 0)

	i := 0
	for _, track := range tracks {
		for _, artist := range track.Artists {
			rows = append(rows, map[string]string{})
			rows[i]["userId"] = userId
			rows[i]["trackId"] = strings.ToLower(strings.ReplaceAll(track.Name, " ", ""))
			for _, artist2 := range track.Artists {
				rows[i]["trackId"] += artist2.Id
			}

			rows[i]["trackName"] = track.Name
			// rows[i]["trackImage"] = track.Album.Image[0].Url
			rows[i]["artistId"] = artist.Id
			rows[i]["artistName"] = artist.Name
			// rows[i]["artistImage"] = artist.Image[0].Url
			rows[i]["albumId"] = track.Album.Id
			rows[i]["albumName"] = track.Album.Name
			// rows[i]["albumImage"] = track.Album.Image[0].Url
			rows[i]["playlistId"] = playlist.Id
			rows[i]["playlistName"] = playlist.Name
			// rows[i]["playlistImage"] = playlist.Image[0].Url
			rows[i]["playlistTotalTracks"] = helpers.ToString(playlist.Tracks.Total)
			rows[i]["playlistLink"] = playlist.ExternalUrls.Spotify
			rows[i]["playlistDescription"] = playlist.Description
			rows[i]["playlistSnapshotId"] = playlist.SnapshotId

			if len(artist.Image) > 0 {
				rows[i]["artistImage"] = artist.Image[0].Url
			} else {
				rows[i]["artistImage"] = ""
			}
			if len(track.Album.Image) > 0 {
				rows[i]["trackImage"] = track.Album.Image[0].Url
				rows[i]["albumImage"] = track.Album.Image[0].Url
			} else {
				rows[i]["trackImage"] = ""
				rows[i]["albumImage"] = ""
			}
			if len(playlist.Image) > 0 {
				rows[i]["playlistImage"] = playlist.Image[0].Url
			} else {
				rows[i]["playlistImage"] = ""
			}
			i += 1
		}
	}
	return rows
}

func (c Client) GetPlaylistsToUpdate(userId string, accessToken string) ([]models.SpotifyPlaylist, error) {
	fetchedPlaylists, err := c.SpotifyClient.GetAllUserPlaylists(userId, accessToken)
	if err != nil {
		return nil, err
	}

	rTransactionResult, err := c.Neo4jReadTransaction(
		"MATCH (user:User)-[:PLAYLIST]->(playlist:Playlist) "+
			"WHERE user.id = $userId "+
			"RETURN playlist.id as playlistId , playlist.total as playlistTotalTracks, playlist.snapshotId as playlistSnapshotId ",
		map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}

	playlistSnapshotIds, err := c.MapPlaylistIdToSnapshotIdFromResult(rTransactionResult)
	if err != nil {
		return nil, err
	}

	if playlistSnapshotIds == nil {
		playlistSnapshotIds = make(map[string]map[string]string)
	}

	for playlistId, playlist := range fetchedPlaylists {
		if _, ok := playlistSnapshotIds[playlistId]; !ok {
			playlistSnapshotIds[playlistId] = make(map[string]string)
			playlistSnapshotIds[playlistId]["totalTracks"] = helpers.ToString(playlist.Tracks.Total)
			playlistSnapshotIds[playlistId]["snapshotId"] = playlist.SnapshotId
			playlistSnapshotIds[playlistId]["exists"] = "1"
		} else if oldSnapShot, ok := playlistSnapshotIds[playlistId]["snapshotId"]; ok {
			if oldSnapShot == playlist.SnapshotId { // playlist has not changed
				delete(playlistSnapshotIds, playlistId)
			} else {
				playlistSnapshotIds[playlistId]["exists"] = "1"
			}
		} else { // new playlist that has not been stored yet
			playlistSnapshotIds[playlistId] = make(map[string]string)
			playlistSnapshotIds[playlistId]["totalTracks"] = helpers.ToString(playlist.Tracks.Total)
			playlistSnapshotIds[playlistId]["snapshotId"] = playlist.SnapshotId
			playlistSnapshotIds[playlistId]["exists"] = "1"
		}
	}

	var playlistsToUpdate []models.SpotifyPlaylist
	var playlistsToRemove []string

	for playlistId, playlistInfo := range playlistSnapshotIds {
		if playlistInfo["exists"] == "1" {
			playlistsToUpdate = append(playlistsToUpdate, fetchedPlaylists[playlistId])
		} else {
			playlistsToRemove = append(playlistsToRemove, playlistId)
		}
	}

	batchDeleteRows := make([]map[string]string, len(playlistsToRemove))
	for i, playlistId := range playlistsToRemove {
		batchDeleteRows[i] = make(map[string]string)
		batchDeleteRows[i]["playlistId"] = playlistId
	}

	_, err = c.Neo4jWriteTransaction(
		"UNWIND $batch as row "+
			"MATCH (playlist:Playlist {id: row.playlistId}) "+
			"DETACH DELETE playlist ",
		map[string]interface{}{"batch": batchDeleteRows},
	)

	if err != nil {
		return nil, err
	}
	return playlistsToUpdate, nil
}

func (c Client) UpdatePlaylist(userId string, accessToken string, playlist models.SpotifyPlaylist) error {
	fetchedTracks, err := c.SpotifyClient.GetAllPlaylistTracks(userId, accessToken, playlist.Id, playlist.Tracks.Total)
	if err != nil {
		return err
	}

	rTransactionResult, err := c.Neo4jReadTransaction(
		"MATCH (playlist:Playlist)-[:PLAYLIST_TRACK]->(track:Track) "+
			"WHERE playlist.id = $playlistId "+
			"RETURN track.id as trackId, track.name as trackName ",
		map[string]interface{}{"playlistId": playlist.Id})
	if err != nil {
		return err
	}

	oldTracks, err := c.MapTrackIdToTrackNameFromResult(rTransactionResult)
	if err != nil {
		log.Println("UPDATE PLAYLIST oldTracks err!: ", err)
		return err
	}

	for trackId := range fetchedTracks {
		if _, ok := oldTracks[trackId]; ok { // track is already stored
			delete(fetchedTracks, trackId)
			delete(oldTracks, trackId)
		}
	}

	batchDeleteRows := make([]map[string]string, len(oldTracks))
	i := 0
	for _, trackId := range oldTracks {
		batchDeleteRows[i] = make(map[string]string)
		batchDeleteRows[i]["playlistId"] = playlist.Id
		batchDeleteRows[i]["trackId"] = trackId
		i += 1
	}

	_, err = c.Neo4jWriteTransaction(
		"UNWIND $batch as row "+
			"MATCH (:Track {id: row.trackId})-[edge1:IN_PLAYLIST]->(:Playlist {id: row.playlistId}) "+
			"MATCH (:Playlist {id: row.playlistId})-[edge2:PLAYLIST_TRACK]->(:Track {id: row.trackId}) "+
			"DELETE edge1 "+
			"DELETE edge2 ",
		map[string]interface{}{"batch": batchDeleteRows},
	)
	if err != nil {
		return err
	}

	// Batch insert
	_, err = c.Neo4jWriteTransaction(
		"UNWIND $batch as row "+
			"MATCH (user:User {id: row.userId}) "+
			"MERGE (track:Track {id: row.trackId}) "+
			"ON CREATE SET track.name = row.trackName, track.image = row.trackImage "+
			"ON MATCH SET track.name = row.trackName, track.image = row.trackImage "+
			"MERGE (album:Album {id: row.albumId}) "+
			"ON CREATE SET album.name = row.albumName, album.image = row.albumImage "+
			"ON MATCH SET album.name = row.albumName, album.image = row.albumImage "+
			"MERGE (artist:Artist {id: row.artistId}) "+
			"ON CREATE SET artist.name = row.artistName, artist.image = row.artistImage "+
			"ON MATCH SET artist.name = row.artistName, artist.image = row.artistImage "+
			"MERGE (playlist:Playlist {id: row.playlistId}) "+
			"ON CREATE SET playlist.name = row.playlistName, playlist.image = row.playlistImage, playlist.snapShotId = row.playlistSnapshotId, playlist.totalTracks = row.playlistTotalTracks,  playlist.link = row.playlistLink, playlist.description = row.playlistDescription "+
			"ON MATCH SET playlist.name = row.playlistName, playlist.image = row.playlistImage, playlist.snapShotId = row.playlistSnapshotId, playlist.totalTracks = row.playlistTotalTracks,  playlist.link = row.playlistLink, playlist.description = row.playlistDescription "+
			"MERGE (user)-[:PLAYLIST]->(playlist) "+
			"MERGE (playlist)-[:PLAYLIST_OWNER]->(user) "+
			"MERGE (track)-[:IN_PLAYLIST]->(playlist) "+
			"MERGE (playlist)-[:PLAYLIST_TRACK]->(track) "+
			"MERGE (artist)-[:IN_TRACK]->(track) "+
			"MERGE (track)-[:TRACK_ARTIST]->(artist) "+
			"MERGE (album)-[:ALBUM_TRACK]->(track) "+
			"MERGE (track)-[:IN_ALBUM]->(album) "+
			"MERGE (artist)-[:ALBUM]->(album) "+
			"MERGE (album)-[:ALBUM_ARTIST]->(artist) ",
		map[string]interface{}{"batch": c.GetBatchInsertRows(userId, fetchedTracks, playlist)},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) UpdatePlaylists(userId string, accessToken string) error {
	playlistsToUpdate, err := c.GetPlaylistsToUpdate(userId, accessToken)
	if err != nil {
		return err
	}

	for _, playlist := range playlistsToUpdate {
		if err != nil {
			return err
		}
		c.UpdatePlaylist(userId, accessToken, playlist)

	}
	return nil
}

func (c Client) UpdateUser(user models.SpotifyUserProfileResponseBody, accessToken string) error {
	userImage := ""
	if len(user.Image) > 0 {
		userImage += user.Image[0].Url
	}
	_, err := c.Neo4jWriteTransaction(
		"MERGE (user:User {id: $id})"+
			"SET user = {id: $id, displayName: $displayName, image: $image, link: $link} "+
			"MERGE (user)-[:GROUP_TYPE]->(:Group {name: $groupName}) ",
		map[string]interface{}{
			"id":          user.Id,
			"displayName": user.DisplayName,
			"image":       userImage,
			"link":        user.ExternalUrls.Spotify,
			"groupName":   "All",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) GetAllUsers(rootUserId string) (map[string]map[string]string, map[string][]string, error) {
	rawUsers, err := c.Neo4jReadTransaction(
		"MATCH (rootUser:User {id: $userId}), (u:User) "+
			"WHERE NOT (rootUser)-[:CONNECTED]-(u:User) "+
			"RETURN  u.id as userId, u.displayName as userDisplayName, u.image as userImage, u.link as userLink ",
		map[string]interface{}{"userId": rootUserId},
	)
	if err != nil {
		return nil, nil, err
	}

	users := make(map[string]map[string]string)    // [userId]{userId, userDisplayName, queryName, userImage}
	queryNametoUserId := make(map[string][]string) // [queryName]{userIds}
	for _, shell := range rawUsers {
		userId, isString1 := shell[0].(string)
		userDisplayName, isString2 := shell[1].(string)
		userImage, isString3 := shell[2].(string)
		userLink, isString4 := shell[3].(string)
		if isString1 && isString2 && isString3 && isString4 {
			queryName := strings.ToLower(strings.ReplaceAll(userDisplayName, " ", ""))
			users[userId] = map[string]string{
				"userId":          userId,
				"userDisplayName": userDisplayName,
				"queryName":       queryName,
				"userImage":       userImage,
				"userLink":        userLink,
			}
			if _, exists := queryNametoUserId[queryName]; !exists {
				queryNametoUserId[queryName] = make([]string, 0)
			}
			queryNametoUserId[queryName] = append(queryNametoUserId[queryName], userId)
		} else {
			log.Println("RECEIVED INCORRECT TYPES FOR userId: ", isString1, " | userDisplayName: ", isString2, " | userImage: ", isString3, " | userLink: ", isString4)
			return nil, nil, fmt.Errorf("RECEIVED INCORRECT TYPES FOR displayName and/or node and/or userImage and/or userLink")
		}
	}
	return users, queryNametoUserId, nil
}

func (c Client) ConnectUsers(rootUserId string, connectUserId string) (bool, error) {
	_, err := c.Neo4jWriteTransaction(
		"MATCH (rootUser:User {id: $rootUserId}), (connectUser:User {id: $connectUserId}) "+
			"MATCH (rootUser)-[:GROUP_TYPE]->(rootUserAllGroup:Group {name: 'All'}) "+
			"MATCH (connectUser)-[:GROUP_TYPE]->(connectUserAllGroup:Group {name: 'All'}) "+
			"MERGE (rootUser)-[:CONNECTED]->(connectUser) "+
			"MERGE (rootUserAllGroup)-[:MEMBER]->(connectUser) "+
			"MERGE (connectUserAllGroup)-[:MEMBER]->(rootUser) ",
		map[string]interface{}{"rootUserId": rootUserId, "connectUserId": connectUserId},
	)
	if err != nil {
		log.Println("CREATING A CONNECTION BETWEEN rootUserId: ", rootUserId, " AND  connectUserId: ", connectUserId, " FAILED")
		return false, err
	}
	return true, nil
}

func (c Client) GetCommonTracks(userId string) (map[string]*models.MusicialSearchTrackResult, map[string][]string, error) {
	rawRecords, err := c.Neo4jReadTransaction(
		"MATCH (rootUser:User {id: $userId})-[:GROUP_TYPE]->(group:Group)-[:MEMBER]->(d1User:User)-[:PLAYLIST]->(d1Playlist:Playlist)-[:PLAYLIST_TRACK]->(d1Track:Track)-[:TRACK_ARTIST]->(d1Artist:Artist) "+
			"RETURN "+
			"group.name as groupName, "+
			"d1User.id as d1UserId, d1User.displayName as d1UserDisplayName,  d1User.image as d1UserImage, d1User.link as d1UserLink, "+
			"d1Playlist.id as d1PlaylistId, d1Playlist.name as d1PlaylistName, d1Playlist.image as d1PlaylistImage, d1Playlist.link as d1PlaylistLink, d1Playlist.description as d1PlaylistDescription, "+
			"d1Track.id as d1TrackId, d1Track.name as d1TrackName, d1Track.image as d1TrackImage, "+
			"d1Artist.id as d1ArtistId, d1Artist.name as d1ArtistName, d1Artist.image as d1ArtistImage ",
		map[string]interface{}{"userId": userId},
	)
	if err != nil {
		return nil, nil, err
	}

	tracks := make(map[string]*models.MusicialSearchTrackResult) // tracks[d1TrackId][d1UserId]{groups = []string, ...(all the other fields as strings)}
	queryNameToTrackId := make(map[string][]string)              // tracks[d1TrackId][d1UserId]{groups = []string, ...(all the other fields as strings)}
	queryNameToTrackIdVis := make(map[string]map[string]struct{})
	artistToSongVis := make(map[string]map[string]struct{})
	songToPlaylistVis := make(map[string]map[string]struct{})
	for _, row := range rawRecords {
		tempRow := make([]string, 16)
		for col := 0; col < 16; col++ {
			item, isString := row[col].(string)
			if !isString {
				log.Panicln("item:", item, " could not be converted to string | getCommonTracks | err: ", err)
				return nil, nil, err
			}
			tempRow[col] = item
		}
		// groupName := tempRow[0]
		d1UserId := tempRow[1]
		d1UserDisplayName := tempRow[2]
		d1UserImage := tempRow[3]
		d1UserLink := tempRow[4]
		d1PlaylistId := tempRow[5]
		d1PlaylistName := tempRow[6]
		d1PlaylistImage := tempRow[7]
		d1PlaylistLink := tempRow[8]
		d1PlaylistDescription := tempRow[9]
		d1TrackId := tempRow[10]
		d1TrackName := tempRow[11]
		d1TrackImage := tempRow[12]
		d1ArtistId := tempRow[13]
		d1ArtistName := tempRow[14]
		d1ArtistImage := tempRow[15]

		if _, exists := tracks[d1TrackId]; !exists {
			tracks[d1TrackId] = &models.MusicialSearchTrackResult{
				TrackId:    d1TrackId,
				TrackName:  d1TrackName,
				TrackImage: d1TrackImage,
				Artists:    make([]models.MusicialSearchTrackArtist, 0),
				Users:      make(map[string]*models.MusicialSearchTrackUser),
			}
			artistToSongVis[d1TrackId] = make(map[string]struct{})

			queryName := strings.ToLower(strings.ReplaceAll(d1TrackName, " ", ""))
			if _, exists := queryNameToTrackId[queryName]; !exists {
				queryNameToTrackId[queryName] = make([]string, 0)
				queryNameToTrackIdVis[queryName] = make(map[string]struct{})
			}
			if _, exists := queryNameToTrackIdVis[queryName][d1TrackId]; !exists {
				queryNameToTrackId[queryName] = append(queryNameToTrackId[queryName], d1TrackId)
				queryNameToTrackIdVis[queryName][d1UserId] = struct{}{}
			}
		}

		if _, exists := artistToSongVis[d1TrackId][d1ArtistId]; !exists {
			tracks[d1TrackId].Artists = append(tracks[d1TrackId].Artists, models.MusicialSearchTrackArtist{
				ArtistId:    d1ArtistId,
				ArtistName:  d1ArtistName,
				ArtistImage: d1ArtistImage,
			})
			artistToSongVis[d1TrackId][d1ArtistId] = struct{}{}
		}

		if _, exists := tracks[d1TrackId].Users[d1UserId]; !exists {
			tracks[d1TrackId].Users[d1UserId] = &models.MusicialSearchTrackUser{
				// Groups:          make(map[string]struct{}),
				UserId:          d1UserId,
				UserDisplayName: d1UserDisplayName,
				UserImage:       d1UserImage,
				UserLink:        d1UserLink,
				Playlists:       make([]models.MusicialSearchTrackUserPlaylist, 0),
			}
			songToPlaylistVis[d1TrackId] = make(map[string]struct{})
		}
		if _, exists := songToPlaylistVis[d1TrackId][d1PlaylistId]; !exists {
			songToPlaylistVis[d1TrackId][d1PlaylistId] = struct{}{}
			tracks[d1TrackId].Users[d1UserId].AddPlaylistToSearchTrackUser(
				models.MusicialSearchTrackUserPlaylist{
					PlaylistId:          d1PlaylistId,
					PlaylistName:        d1PlaylistName,
					PlaylistImage:       d1PlaylistImage,
					PlaylistDescription: d1PlaylistDescription,
					PlaylistLink:        d1PlaylistLink,
				},
			)
		}
		// tracks[d1TrackId].Users[d1UserId].Groups[groupName] = struct{}{}
	}
	return tracks, queryNameToTrackId, nil
}

func (c Client) GetCommonArtists(userId string) (map[string]*models.MusicialSearchArtistResult, map[string][]string, error) {
	rawRecords, err := c.Neo4jReadTransaction(
		"MATCH (rootUser:User {id: $userId})-[:GROUP_TYPE]->(group:Group)-[:MEMBER]->(d1User:User)-[:PLAYLIST]->(d1Playlist:Playlist)-[:PLAYLIST_TRACK]->(d1Track:Track)-[:TRACK_ARTIST]->(d1Artist:Artist) "+
			"RETURN "+
			"group.name as groupName, "+
			"d1User.id as d1UserId, d1User.displayName as d1UserDisplayName,  d1User.image as d1UserImage, d1User.link as d1UserLink, "+
			"d1Playlist.id as d1PlaylistId, d1Playlist.name as d1PlaylistName, d1Playlist.image as d1PlaylistImage, d1Playlist.link as d1PlaylistLink, d1Playlist.description as d1PlaylistDescription, "+
			"d1Track.id as d1TrackId, d1Track.name as d1TrackName, d1Track.image as d1TrackImage, "+
			"d1Artist.id as d1ArtistId, d1Artist.name as d1ArtistName ",
		map[string]interface{}{"userId": userId},
	)
	if err != nil {
		return nil, nil, err
	}

	artists := make(map[string]*models.MusicialSearchArtistResult) // tracks[d1TrackId][d1UserId]{groups = []string, ...(all the other fields as strings)}
	queryNameToArtistId := make(map[string][]string)               // tracks[d1TrackId][d1UserId]{groups = []string, ...(all the other fields as strings)}
	queryNameToArtistIdVis := make(map[string]map[string]struct{})
	songToPlaylistVis := make(map[string]map[string]struct{})
	for _, row := range rawRecords {
		tempRow := make([]string, 15)
		for col := 0; col < 15; col++ {
			item, isString := row[col].(string)
			if !isString {
				log.Panicln("item:", item, " could not be converted to string | getCommonArtists | err: ", err)
				return nil, nil, err
			}
			tempRow[col] = item
		}
		// groupName := tempRow[0]
		d1UserId := tempRow[1]
		d1UserDisplayName := tempRow[2]
		d1UserImage := tempRow[3]
		d1UserLink := tempRow[4]
		d1PlaylistId := tempRow[5]
		d1PlaylistName := tempRow[6]
		d1PlaylistImage := tempRow[7]
		d1PlaylistLink := tempRow[8]
		d1PlaylistDescription := tempRow[9]
		d1TrackId := tempRow[10]
		d1TrackName := tempRow[11]
		d1TrackImage := tempRow[12]
		d1ArtistId := tempRow[13]
		d1ArtistName := tempRow[14]

		if _, exists := artists[d1ArtistId]; !exists {
			artists[d1ArtistId] = &models.MusicialSearchArtistResult{
				ArtistId:   d1ArtistId,
				ArtistName: d1ArtistName,
				Users:      make(map[string]*models.MusicialSearchArtistUser),
			}
			// artistToPlaylistVis[d1PlaylistId] = make(map[string]struct{})

			queryName := strings.ToLower(strings.ReplaceAll(d1ArtistName, " ", ""))
			if _, exists := queryNameToArtistId[queryName]; !exists {
				queryNameToArtistId[queryName] = make([]string, 0)
				queryNameToArtistIdVis[queryName] = make(map[string]struct{})
			}
			if _, exists := queryNameToArtistIdVis[queryName][d1ArtistId]; !exists {
				queryNameToArtistId[queryName] = append(queryNameToArtistId[queryName], d1ArtistId)
				queryNameToArtistIdVis[queryName][d1ArtistId] = struct{}{}
			}
		}

		if _, exists := artists[d1ArtistId].Users[d1UserId]; !exists {
			artists[d1ArtistId].Users[d1UserId] = &models.MusicialSearchArtistUser{
				UserId:          d1UserId,
				UserDisplayName: d1UserDisplayName,
				UserImage:       d1UserImage,
				UserLink:        d1UserLink,
				Tracks:          make(map[string]*models.MusicialSearchArtistUserTrack),
			}
		}

		if _, exists := artists[d1ArtistId].Users[d1UserId].Tracks[d1TrackId]; !exists {
			artists[d1ArtistId].Users[d1UserId].Tracks[d1TrackId] = &models.MusicialSearchArtistUserTrack{
				TrackId:    d1TrackId,
				TrackName:  d1TrackName,
				TrackImage: d1TrackImage,
				Playlists:  make([]models.MusicialSearchArtistTrackPlaylist, 0),
			}
			songToPlaylistVis[d1TrackId] = make(map[string]struct{})
		}

		if _, exists := songToPlaylistVis[d1TrackId][d1PlaylistId]; !exists {
			songToPlaylistVis[d1TrackId][d1PlaylistId] = struct{}{}
			artists[d1ArtistId].Users[d1UserId].Tracks[d1TrackId].AddPlaylistToSearchArtistUserTrack(
				models.MusicialSearchArtistTrackPlaylist{
					PlaylistId:          d1PlaylistId,
					PlaylistName:        d1PlaylistName,
					PlaylistImage:       d1PlaylistImage,
					PlaylistDescription: d1PlaylistDescription,
					PlaylistLink:        d1PlaylistLink,
				},
			)
		}
		// artists[d1TrackId].Users[d1UserId].Groups[groupName] = struct{}{}
	}
	return artists, queryNameToArtistId, nil
}
