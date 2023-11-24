# Spotify Playlist Exporter

A Go script that uses `goroutines` and `waitgroups` that converts a Spotify playlist into a .csv file.

## Inputs:

- `playlistUrls`: A slice of strings. These strings are Spotify playlist urls.

## Outputs:

- `./exports/<TIMESTAMP>/<PLAYLIST_NAME>.csv`: .csv files that contain important information about the playlist's tracks such as:
  - Name
  - Artists
  - Album Name
  - Duration
