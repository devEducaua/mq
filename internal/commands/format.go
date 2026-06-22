package commands

import (
	"fmt"
	"mq/internal/mpd"
	"time"
	"unicode/utf8"
)

func PrintFormattedQueue(queue []mpd.Song) error {
	maxTitleLength := 0;
	maxAlbumLength := 0;
	maxArtistLength := 0;

	current, err := mpd.GetCurrentSong();
	if err != nil {
		return err;
	}

	for _,s := range queue {
		maxTitleLength = max(maxTitleLength, utf8.RuneCountInString(s.Title));
		maxAlbumLength = max(maxAlbumLength, utf8.RuneCountInString(s.Album));
		maxArtistLength = max(maxArtistLength, utf8.RuneCountInString(s.Artist));
	}

	if len(queue) == 1 {
		return nil;
	}

	for _,s := range queue {
		var marker = " ";
		if s.File == current.File {
			marker = "*";
		}
		fmt.Printf("%-4v%1s %v %-*s - %-*s - %-*s\n", s.Pos, marker, FormatDuration(float64(s.Time)), maxTitleLength, s.Title, maxAlbumLength, s.Album, maxArtistLength, s.Artist);
	}
	return nil;
}


func PrintFormattedStatus(plainResponse string) error {
	current, err := mpd.GetCurrentSong();
	if err != nil {
		return err;
	}
	status, err := mpd.ParseStatusResponse(plainResponse);
	if err != nil {
		return err;
	}

	artist := current.Artist;
	if artist == "" {
		artist = "no artist"
	}
	title := current.Title;
	if title == "" {
		title = "no title";
	}

	fmt.Printf("%v - %v\n", artist, title);
	fmt.Printf("%v/%v - %v/%v - state: %v\n", status.Song, status.Playlistlength, FormatDuration(status.Elapsed), FormatDuration(status.Duration), status.State);
	fmt.Printf("repeat: %v, random: %v, single: %v, consume: %v\n", status.Repeat, status.Random, status.Single, status.Consume)

	return nil;
}

func FormatDuration(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second));

	hours := int(d.Hours());
	minutes := int(d.Minutes()) % 60;
	secs := int(d.Seconds()) % 60;

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs);
	}

	return fmt.Sprintf("%d:%02d", minutes, secs);
}

