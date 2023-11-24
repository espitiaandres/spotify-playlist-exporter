package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"

	"spotify-playlist-exporter/helpers/csvfiles"
	"spotify-playlist-exporter/helpers/spotifyclient"
	"spotify-playlist-exporter/helpers/timer"
)

/**
TODO:
- complete readme.md
*/

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
		"https://open.spotify.com/playlist/37i9dQZF1DWV7EzJMK2FUI?si=7e466243d8e84189",
		"https://open.spotify.com/playlist/4uClt6zoLaRF0vHPEwWChR?si=7899b77a5f594b3f",
		"https://open.spotify.com/playlist/6A38ofY5uoJpFLCwE2T8OC?si=7acb4e42de0b409a",
	}

	currentTime := time.Now()
	currentTimeFormatted := strings.ReplaceAll(currentTime.Format("01-02-2006 15:04:05"), " ", "_")

	var wg sync.WaitGroup

	for _, element := range playlistUrls {
		wg.Add(1)
		playlistQueryString := strings.SplitAfter(element, "https://open.spotify.com/playlist/")[1]
		playlistIdRaw := strings.SplitAfter(playlistQueryString, "?si=")[0]
		playlistId := strings.ReplaceAll(playlistIdRaw, "?si=", "")

		go func() {
			defer wg.Done()
			client := spotify.Authenticator{}.NewClient(accessToken)
			playlist := spotifyclient.FetchSpotifyClient(accessToken, playlistId)
			csvfiles.CreateCsv(client, accessToken, playlist, currentTimeFormatted)
		}()
	}

	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	log.Println("Waiting for goroutines to complete...")
	wg.Wait()
	log.Printf("Completed converting %v playlists to .csv", len(playlistUrls))
}
