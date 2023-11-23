package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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
		"https://open.spotify.com/playlist/37i9dQZF1DWV7EzJMK2FUI?si=7e466243d8e84189",
	}

	var wg sync.WaitGroup

	for _, element := range playlistUrls {
		wg.Add(1)
		playlistQueryString := strings.SplitAfter(element, "https://open.spotify.com/playlist/")[1]
		playlistId := strings.SplitAfter(playlistQueryString, "?si=")[0]

		go func() {
			defer wg.Done()
			fetchSpotifyClient(accessToken, playlistId)
		}()
	}

	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	log.Println("Waiting for goroutines to complete...")
	wg.Wait()
	log.Printf("Completed converting %v playlists to .csv", len(playlistUrls))
}

func fetchSpotifyClient(accessToken *oauth2.Token, playlistId string) {
	log.Println("playlist id:", playlistId)

	// client := spotify.Authenticator{}.NewClient(accessToken)

	spotifyPlaylistId := spotify.ID(playlistId)

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlist, err := client.GetPlaylist(spotifyPlaylistId)

	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	createCsv(accessToken, playlist)
}

func createCsv(accessToken *oauth2.Token, playlist *spotify.FullPlaylist) {
	playlistName := playlist.Name
	playlistNameFormatted := strings.ReplaceAll(playlistName, " ", "_")

	currentTime := time.Now()
	currentTimeFormatted := strings.ReplaceAll(currentTime.Format("01-02-2006 15:04:05"), " ", "_")

	fileName := fmt.Sprintf("./exports/%s_%s.csv", playlistNameFormatted, currentTimeFormatted)

	tracks := playlist.Tracks.Tracks
	var allTracks [][]string
	titleRow := []string{
		"Track",
		"Artists",
		"Album",
		"Duration",
	}

	allTracks = append(allTracks, titleRow)

	for _, element := range tracks {
		var artists []string

		for _, artist := range element.Track.Artists {
			artists = append(artists, artist.Name)
		}

		trackName := element.Track.Name
		artistsString := strings.Join(artists, ", ")
		albumName := element.Track.Album.Name
		trackDuration := element.Track.TimeDuration().String()

		row := []string{
			trackName,
			artistsString,
			albumName,
			trackDuration,
		}

		allTracks = append(allTracks, row)
	}

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	defer file.Close()

	wr := csv.NewWriter(file)
	wr.WriteAll(allTracks)
}
