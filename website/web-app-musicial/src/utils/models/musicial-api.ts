export interface Neo4jAllUserWrapper {
  users: Map<string, UserWithQueryName>;
  queryNameToUserId: Map<string, string[]>;
}

export interface UserWithQueryName {
  userId: string;
  userDisplayName: string;
  userImage: string;
  userLink: string;
  queryName: string;
}

export interface Neo4jAllTrackWrapper {
  tracks: Map<string, TrackResult>;
  queryNameToTrackId: Map<string, string[]>;
}

export interface TrackResult {
  trackId: string;
  trackName: string;
  trackImage: string;
  artists: Array<TrackResultArtist>;
  users: Map<string, TrackResultUser>;
}

export interface TrackResultUser {
  // groups: Set<string>;
  userId: string;
  userDisplayName: string;
  userImage: string;
  userLink: string;
  playlists: Array<TrackResultPlaylist>;
}

export interface TrackResultPlaylist {
  playlistId: string;
  playlistName: string;
  playlistImage: string;
  playlistDescription: string;
  playlistLink: string;
}

export interface TrackResultArtist {
  artistId: string;
  artistName: string;
  artistImage: string;
}

export interface TrackSuggestion {
  trackId: string;
  trackName: string;
  trackImage: string;
  trackArtists: string;
}

export interface Neo4jAllArtistWrapper {
  artists: Map<string, ArtistResult>;
  queryNameToArtistId: Map<string, string[]>;
}

export interface ArtistResult {
  artistId: string;
  artistName: string;
  users: Map<string, ArtistResultUser>;
}

export interface ArtistResultUser {
  userId: string;
  userDisplayName: string;
  userImage: string;
  userLink: string;
  tracks: Map<string, ArtistResultTrack>;
}

export interface ArtistResultTrackPlaylist {
  playlistId: string;
  playlistName: string;
  playlistImage: string;
  playlistDescription: string;
  playlistLink: string;
}

export interface ArtistResultTrack {
  trackId: string;
  trackName: string;
  trackImage: string;
  playlists: ArtistResultTrackPlaylist[];
}

export interface ArtistSuggestion {
  artistId: string;
  artistName: string;
}

export interface UserProperties {
  image: string;
  id: string;
  displayName: string;
}

export interface TrackProperties {
  name: string;
  image: string;
  id: string;
}

export interface ArtistProperties {
  name: string;
  image: string;
  id: string;
}

export interface AlbumProperties {
  name: string;
  image: string;
  id: string;
}

export interface PlaylistProperties {
  name: string;
  image: string;
  id: string;
}
