package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

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
- move functions out
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
			csvfiles.CreateCsv(client, accessToken, playlist)
		}()
	}

	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	log.Println("Waiting for goroutines to complete...")
	wg.Wait()
	log.Printf("Completed converting %v playlists to .csv", len(playlistUrls))
}

// func createCsv(accessToken *oauth2.Token, playlist *spotify.FullPlaylist) {
// 	tracks := playlist.Tracks.Tracks
// 	var allTracks [][]string
// 	titleRow := []string{
// 		"Track",
// 		"Artists",
// 		"Album",
// 		"Duration",
// 	}

// 	allTracks = append(allTracks, titleRow)

// 	for _, element := range tracks {
// 		var artists []string

// 		for _, artist := range element.Track.Artists {
// 			artists = append(artists, artist.Name)
// 		}

// 		trackName := element.Track.Name
// 		artistsString := strings.Join(artists, ", ")
// 		albumName := element.Track.Album.Name
// 		trackDuration := element.Track.TimeDuration().String()

// 		row := []string{
// 			trackName,
// 			artistsString,
// 			albumName,
// 			trackDuration,
// 		}

// 		allTracks = append(allTracks, row)
// 	}

// 	fileName := createCsvFileName(playlist)

// 	file, err := os.Create(fileName)

// 	if err != nil {
// 		log.Fatalln("failed to open file", err)
// 	}

// 	defer file.Close()

// 	wr := csv.NewWriter(file)
// 	wr.WriteAll(allTracks)
// }

// func createCsvFileName(playlist *spotify.FullPlaylist) string {
// 	playlistName := playlist.Name
// 	playlistNameFormatted := strings.ReplaceAll(playlistName, " ", "_")

// 	currentTime := time.Now()
// 	currentTimeFormatted := strings.ReplaceAll(currentTime.Format("01-02-2006 15:04:05"), " ", "_")

// 	if err := os.MkdirAll(fmt.Sprintf("./exports/%s/", currentTimeFormatted), os.ModePerm); err != nil {
// 		log.Fatal(err)
// 	}

// 	fileName := fmt.Sprintf("./exports/%s/%s.csv", currentTimeFormatted, playlistNameFormatted)

// 	return fileName
// }
