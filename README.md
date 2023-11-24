# Spotify Playlist Exporter

``spotify-playlist-exporter` is a Go script that uses `goroutines` and `waitgroups` that converts a Spotify playlist into a .csv file.

## Table of contents:

- [Usage](#usage)
  - [Inputs](#inputs)
  - [Outputs](#outputs)
- [Documentation](#documentation)
- [Bugs/Requests](#bugs_requests)
- [License](#license)

<a name="usage"/>

# Usage

Below are the inputs/outputs to this script:

<a name="inputs"/>

## Inputs:

- `playlistUrls`: A slice of strings. These strings are Spotify playlist urls.

> Example: https://open.spotify.com/playlist/37i9dQZEVXbNG2KDcFcKOF?si=be29ccfd84a347b8

<a name="outputs"/>

## Outputs:

- `./exports/<TIMESTAMP>/<PLAYLIST_NAME>.csv`: .csv files that contain important information about the playlist's tracks such as:
  - Name
  - Artists
  - Album Name
  - Duration

> Example:

<h1 align="center">
<img
    src="/assets/example_csv_file.jpg"
    alt="Example CSV file produced"
    title="Example CSV file produced"
    width="200"
/>
</h1>

<a name="documentation"/>

# Documentation

Documentation of `penguin-py` can be found here: https://github.com/espitiaandres/spotify-playlist-exporter/blob/main/README.md

<a name="bugs_requests"/>

# Bugs/Requests

If you find any bugs or have any suggestions to `spotify-playlist-exporter`, submit them in the issues tab in the Github repo. This can be found here: https://github.com/espitiaandres/spotify-playlist-exporter/issues

<a name="license"/>

# License

Distributed under the terms of the MIT license, `spotify-playlist-exporter` is free and open source software.
