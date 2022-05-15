package models

type SpotifyTrack struct {
	Album            SpotifyAlbum        `json:"album,omitempty"`
	Artists          []SpotifyArtist     `json:"artists,omitempty"`
	AvailableMarkets []string            `json:"available_markets,omitempty"`
	DiscNumber       int                 `json:"disc_number,omitempty"`
	DurationMs       int                 `json:"duration_ms,omitempty"`
	Explicit         bool                `json:"explicit,omitempty"`
	ExternalIds      SpotifyExternalIds  `json:"external_ids,omitempty"`
	ExternalUrls     SpotifyExternalUrls `json:"external_urls,omitempty"`
	Href             string              `json:"href,omitempty"`
	Id               string              `json:"id,omitempty"`
	Name             string              `json:"name,omitempty"`
	Popularity       string              `json:"popularity,omitempty"`
	TrackNumber      int                 `json:"track_number,omitempty"`
	Type             string              `json:"type,omitempty"`
	Uri              string              `json:"uri,omitempty"`
	IsLocal          bool                `json:"is_local,omitempty"`
}

type SpotifyPlaylistTrack struct {
	AddedAt       string                `json:"added_at,omitempty"`
	AddedBy       string                `json:"added_by,omitempty"`
	IsLocal       bool                  `json:"is_local,omitempty"`
	Track         SpotifyTrack          `json:"track,omitempty"`
	VideoThumbnal SpotifyVideoThumbNail `json:"video_thumbnail,omitempty"`
	// PrimaryColor  string                `json:"primary_color,omitempty"`
}

type SpotifyTracks struct {
	Href     string         `json:"href,omitempty"`
	Items    []SpotifyTrack `json:"items,omitempty"`
	Limit    int            `json:"limit,omitempty"`
	Next     string         `json:"next,omitempty"`
	Offset   int            `json:"ofset,omitempty"`
	Previous string         `json:"previous,omitempty"`
	Total    int            `json:"total,omitempty"`
}

type SpotifyArtist struct {
	ExternalUrls SpotifyExternalUrls `json:"external_urls,omitempty"`
	Followers    SpotifyFollowers    `json:"followers,omitempty"`
	Genres       []string            `json:"genres,omitempty"`
	Href         string              `json:"href,omitempty"`
	Id           string              `json:"id,omitempty"`
	Image        []SpotifyImage      `json:"images,omitempty"`
	Name         string              `json:"name,omitempty"`
	Popularity   int                 `json:"popularity,omitempty"`
	Type         string              `json:"type,omitempty"`
	Uri          string              `json:"uri,omitempty"`
}

type SpotifyAlbum struct {
	AlbumType            string              `json:"album,omitempty"`
	TotalTracks          int                 `json:"total_tracks,omitempty"`
	AvailableMarkets     []string            `json:"available_markets,omitempty"`
	ExternalUrls         SpotifyExternalUrls `json:"external_urls,omitempty"`
	Href                 string              `json:"href,omitempty"`
	Id                   string              `json:"id,omitempty"`
	Image                []SpotifyImage      `json:"images,omitempty"`
	Name                 string              `json:"name,omitempty"`
	ReleaseDate          string              `json:"release_date,omitempty"`
	ReleaseDatePrecision string              `json:"release_date_precision,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Uri                  string              `json:"uri,omitempty"`
	Artists              []SpotifyArtist     `json:"artists,omitempty"`
	Tracks               SpotifyTracks       `json:"tracks,omitempty"`
}

type SpotifyPlaylist struct {
	Collaborative bool                 `json:"collaborative,omitempty"`
	Description   string               `json:"description,omitempty"`
	ExternalUrls  SpotifyExternalUrls  `json:"external_urls,omitempty"`
	Followers     SpotifyFollowers     `json:"followers,omitempty"`
	Href          string               `json:"href,omitempty"`
	Id            string               `json:"id,omitempty"`
	Image         []SpotifyImage       `json:"images,omitempty"`
	Name          string               `json:"name,omitempty"`
	Owner         SpotifyPlaylistOwner `json:"owner,omitempty"`
	Public        bool                 `json:"public,omitempty"`
	SnapshotId    string               `json:"snapshot_id,omitempty"`
	Tracks        SpotifyTracks        `json:"tracks,omitempty"`
	Type          string               `json:"type,omitempty"`
	Uri           string               `json:"uri,omitempty"`
}

type SpotifyPlaylistOwner struct {
	ExternalUrls SpotifyExternalUrls `json:"external_urls,omitempty"`
	Followers    SpotifyFollowers    `json:"followers,omitempty"`
	Href         string              `json:"href,omitempty"`
	Id           string              `json:"id,omitempty"`
	Type         string              `json:"type,omitempty"`
	Uri          string              `json:"uri,omitempty"`
	DisplayName  string              `json:"display_name,omitempty"`
}

type SpotifyExplicitContent struct {
	FilterEnabled bool `json:"filter_enabled,omitempty"`
	FilterLocked  bool `json:"filter_locked,omitempty"`
}

type SpotifyVideoThumbNail struct {
	Url string `json:"url,omitempty"`
}

type SpotifyTrackAddedBy struct {
	ExternalUrls SpotifyExternalUrls `json:"external_urls,omitempty"`
	Href         string              `json:"href,omitempty"`
	Id           string              `json:"id,omitempty"`
	Type         string              `json:"type,omitempty"`
	Uri          string              `json:"uri,omitempty"`
}

type SpotifyImage struct {
	Url    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type SpotifyFollowers struct {
	Href  string `json:"href,omitempty"`
	Total int    `json:"total,omitempty"`
}

type SpotifyExternalUrls struct {
	Spotify string `json:"spotify,omitempty"`
}

type SpotifyExternalIds struct {
	Isrc string `json:"isrc,omitempty"`
	Ean  string `json:"ean,omitempty"`
	Upc  string `json:"upc,omitempty"`
}

type SpotifyTokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyTokenResponseBodyWithUserId struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"userId,omitempty"`
}

type SpotifyUserProfileResponseBody struct {
	Country         string                 `json:"country,omitempty"`
	DisplayName     string                 `json:"display_name,omitempty"`
	Email           string                 `json:"email,omitempty"`
	ExplicitContent SpotifyExplicitContent `json:"explicit_content,omitempty"`
	ExternalUrls    SpotifyExternalUrls    `json:"external_urls,omitempty"`
	Followers       SpotifyFollowers       `json:"followers,omitempty"`
	Href            string                 `json:"href,omitempty"`
	Id              string                 `json:"id,omitempty"`
	Image           []SpotifyImage         `json:"images,omitempty"`
	Product         string                 `json:"product,omitempty"`
	Type            string                 `json:"type,omitempty"`
	Uri             string                 `json:"uri,omitempty"`
}

type SpotifyPlaylistTracksResponseBody struct {
	Href     string                 `json:"href,omitempty"`
	Items    []SpotifyPlaylistTrack `json:"items,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Next     string                 `json:"next,omitempty"`
	Offset   int                    `json:"offset,omitempty"`
	Previous string                 `json:"previous,omitempty"`
	Total    int                    `json:"total,omitempty"`
}

type SpotifyUserPlaylistsResponseBody struct {
	Href     string            `json:"href,omitempty"`
	Items    []SpotifyPlaylist `json:"items,omitempty"`
	Limit    int               `json:"limit,omitempty"`
	Next     string            `json:"next,omitempty"`
	Offset   int               `json:"offset,omitempty"`
	Previous string            `json:"previous,omitempty"`
	Total    int               `json:"total,omitempty"`
}
