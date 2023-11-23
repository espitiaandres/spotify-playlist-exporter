package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"spotify-playlist-exporter/helpers/timer"
)

func main() {
	defer timer.FuncTimer("main")()

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	authConfig := &clientcredentials.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	contextBackground := context.Background()
	accessToken, err := authConfig.Token(contextBackground)

	playlistUrls := []string{
		"https://open.spotify.com/playlist/37i9dQZF1DWV7EzJMK2FUI?si=799ca1af48274f9d",
		"https://open.spotify.com/playlist/37i9dQZF1DX0jgyAiPl8Af?si=edf2c8829a4c4071",
	}

	var wg sync.WaitGroup

	for _, element := range playlistUrls {
		wg.Add(1)
		playlistQueryString := strings.SplitAfter(element, "https://open.spotify.com/playlist/")[1]
		playlistId := strings.SplitAfter(playlistQueryString, "?si=")[0]

		go func() {
			defer wg.Done()
			fetchPlaylist(accessToken, playlistId)
		}()
	}

	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	log.Println("Waiting for goroutines to complete...")
	wg.Wait()
	log.Printf("Completed converting %v playlists to .csv", len(playlistUrls))
}

func fetchPlaylist(accessToken *oauth2.Token, playlistId string) {
	log.Println("playlist id:", playlistId)

	client := spotify.Authenticator{}.NewClient(accessToken)

	spotifyPlaylistId := spotify.ID(playlistId)

	playlist, err := client.GetPlaylist(spotifyPlaylistId)

	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	createCsv(playlist)
}

func createCsv(playlist *spotify.FullPlaylist) {
	playlistName := playlist.Name

	playlistNameFormatted := strings.ReplaceAll(playlistName, " ", "_")

	fileName := fmt.Sprintf("./exports/%s.csv", playlistNameFormatted)

	log.Println("playlist name:", playlist.Name)
	log.Println("filename:", fileName)
}
