package models

type MusicialUser struct {
	Id          string
	DisplayName string
	Image       string
}

type MusicialTrack struct {
	Id    string
	Name  string
	Image string
}

type MusicialAlbum struct {
	Id    string
	Name  string
	Image string
}

type MusicialArtist struct {
	Id    string
	Name  string
	Image string
}

type MusicialPlaylist struct {
	Id          string
	Name        string
	TotalTracks string
	Image       string
	Description string
	SnapshotId  string
}

type MusicialSearchTrackResult struct {
	TrackId    string                              `json:"trackId"`
	TrackName  string                              `json:"trackName"`
	TrackImage string                              `json:"trackImage"`
	Artists    []MusicialSearchTrackArtist         `json:"artists"`
	Users      map[string]*MusicialSearchTrackUser `json:"users"`
}

type MusicialSearchTrackUser struct {
	UserId          string                            `json:"userId"`
	UserDisplayName string                            `json:"userDisplayName"`
	UserImage       string                            `json:"userImage"`
	UserLink        string                            `json:"userLink"`
	Playlists       []MusicialSearchTrackUserPlaylist `json:"playlists"`
}

type MusicialSearchTrackArtist struct {
	ArtistId    string `json:"artistId"`
	ArtistName  string `json:"artistName"`
	ArtistImage string `json:"artistImage"`
}

type MusicialSearchTrackUserPlaylist struct {
	PlaylistId          string `json:"playlistId"`
	PlaylistName        string `json:"playlistName"`
	PlaylistImage       string `json:"playlistImage"`
	PlaylistDescription string `json:"playlistDescription"`
	PlaylistLink        string `json:"playlistLink"`
}

func (u *MusicialSearchTrackUser) AddPlaylistToSearchTrackUser(playlist MusicialSearchTrackUserPlaylist) {
	u.Playlists = append(u.Playlists, playlist)
}

type MusicialSearchArtistResult struct {
	ArtistId   string `json:"artistId"`
	ArtistName string `json:"artistName"`
	// ArtistImage string `json:"artistImage"`
	Users map[string]*MusicialSearchArtistUser `json:"users"`
}

type MusicialSearchArtistUser struct {
	UserId          string                                    `json:"userId"`
	UserDisplayName string                                    `json:"userDisplayName"`
	UserImage       string                                    `json:"userImage"`
	UserLink        string                                    `json:"userLink"`
	Tracks          map[string]*MusicialSearchArtistUserTrack `json:"tracks"`
}

type MusicialSearchArtistUserTrack struct {
	TrackId    string                              `json:"trackId"`
	TrackName  string                              `json:"trackName"`
	TrackImage string                              `json:"trackImage"`
	Playlists  []MusicialSearchArtistTrackPlaylist `json:"playlists"`
}

type MusicialSearchArtistTrackPlaylist struct {
	PlaylistId          string `json:"playlistId"`
	PlaylistName        string `json:"playlistName"`
	PlaylistImage       string `json:"playlistImage"`
	PlaylistDescription string `json:"playlistDescription"`
	PlaylistLink        string `json:"playlistLink"`
}

func (p *MusicialSearchArtistUserTrack) AddPlaylistToSearchArtistUserTrack(playlist MusicialSearchArtistTrackPlaylist) {
	p.Playlists = append(p.Playlists, playlist)
}

type LambdaPayload struct {
	// Method string `json:"httpMethod"`
	Body string `json:"body"`
}
