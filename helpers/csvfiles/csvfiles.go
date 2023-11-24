package csvfiles

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func CreateCsv(client spotify.Client, accessToken *oauth2.Token, playlist *spotify.FullPlaylist, currentTimeFormatted string) {
	tracks, err := client.GetPlaylistTracks(playlist.ID)

	playlistTracks := []spotify.PlaylistTrack{}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Playlist has %d total tracks", tracks.Total)

	for page := 1; ; page++ {
		playlistTracks = append(playlistTracks, tracks.Tracks...)

		err = client.NextPage(tracks)

		if err == spotify.ErrNoMorePages {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	var allTracks [][]string
	titleRow := []string{
		"Track",
		"Artists",
		"Album",
		"Duration",
	}

	allTracks = append(allTracks, titleRow)

	for _, element := range playlistTracks {
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

	fileName := createCsvFileName(playlist, currentTimeFormatted)

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	defer file.Close()

	wr := csv.NewWriter(file)
	wr.WriteAll(allTracks)
}

func createCsvFileName(playlist *spotify.FullPlaylist, currentTimeFormatted string) string {
	playlistName := playlist.Name
	playlistNameFormatted := strings.ReplaceAll(playlistName, " ", "_")

	directory := fmt.Sprintf("./exports/%s/", currentTimeFormatted)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf("./exports/%s/%s.csv", currentTimeFormatted, playlistNameFormatted)

	return fileName
}
