package spotifyclient

import (
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func FetchSpotifyClient(accessToken *oauth2.Token, playlistId string) *spotify.FullPlaylist {
	spotifyPlaylistId := spotify.ID(playlistId)

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlist, err := client.GetPlaylist(spotifyPlaylistId)

	if err != nil {
		log.Fatalf("error retrieve playlist with id: %s. Error: %v", spotifyPlaylistId, err)
	}

	return playlist
}
