package csvfiles

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func CreateCsv(client spotify.Client, accessToken *oauth2.Token, playlist *spotify.FullPlaylist) {
	var limit = 100

	playlistTracks, err := client.GetPlaylistTracksOpt(playlist.ID, &spotify.Options{
		Limit: &limit,
	}, "limit")

	if err != nil {
		log.Fatalln("No playlist tracks", err)
	}

	tracks := playlist.Tracks.Tracks

	search, err := client.Search(playlist.Tracks.Next, spotify.SearchTypePlaylist)

	if err != nil {
		log.Fatalln("Fail to search", err)
	}

	// TODO: loop here.
	// If playlist.Tracks.next is not nil, keep adding them to the tracks array
	// client.NextTrackResults(search)

	// tracks = append(tracks, search.Tracks.Tracks)

	fmt.Println(playlist.Name, playlist.Tracks.Total, playlistTracks.Limit, search.Tracks)

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

	fileName := createCsvFileName(playlist)

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	defer file.Close()

	wr := csv.NewWriter(file)
	wr.WriteAll(allTracks)
}

func createCsvFileName(playlist *spotify.FullPlaylist) string {
	playlistName := playlist.Name
	playlistNameFormatted := strings.ReplaceAll(playlistName, " ", "_")

	currentTime := time.Now()
	currentTimeFormatted := strings.ReplaceAll(currentTime.Format("01-02-2006 15:04:05"), " ", "_")

	if err := os.MkdirAll(fmt.Sprintf("./exports/%s/", currentTimeFormatted), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf("./exports/%s/%s.csv", currentTimeFormatted, playlistNameFormatted)

	return fileName
}
