package constants

const BaseApiUrl = "https://api.musicial.net/api/"
const ApiProcessLoginUrl = "https://api.musicial.net/api/login/process"

// Spotify Constants
const BaseSpotifyAuthUrl = "https://accounts.spotify.com/authorize?"
const SpotifyTokenURL = "https://accounts.spotify.com/api/token"
const SpotifyUserURL = "https://api.spotify.com/v1/me"
const StateKey = "spotify_auth_state"

const SpotifyUserPlaylistsURL = "https://api.spotify.com/v1/me/playlists"
const SpotifyPlaylistTracks = "https://api.spotify.com/v1/playlists/%s/tracks"
const SpotifyGetPlaylistItemsFields = "items(track(album(id,name,artists(id,name,images),images),artists(id,name,images),id,name))"

// Cookie Cosntants
const InviterCookieKey = "inviter"
