package mpd

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Song struct {
	File string
	Format string
	Artist string
	Title string
	Album string
	Date int
	Time int
	Duration int
	Pos int
	Id int
}

type Status struct {
	Volume int
	Repeat int
	Random int
	Single string
	Consume string
	Playlist int
	Playlistlength int
	State string
	Song int
	SongId int
	Duration float64
	Elapsed float64
	NextSong int
	NextSongId int
}

func ParseBinaryResponse(reader *bufio.Reader) ([]byte, int, error) {

	// size:
	line, err := reader.ReadString('\n');
	if err != nil {
		return nil, 0, err;
	}
	parts := strings.SplitN(line, ":", 2);
	value := strings.TrimSpace(parts[1]);

	size, err := strconv.Atoi(value);
	if err != nil {
		return nil, 0, err;
	}

	// type:
	_, err = reader.ReadString('\n');
	if err != nil {
		return nil, 0, err;
	}

	// binary:
	line, err = reader.ReadString('\n');
	if err != nil {
		return nil, 0, err;
	}
	parts = strings.SplitN(line, ":", 2);
	value = strings.TrimSpace(parts[1]);

	binary, err := strconv.Atoi(value);
	if err != nil {
		return nil, 0, err;
	}

	// binary data
	buf := make([]byte, binary);
	_, err = io.ReadFull(reader, buf);
	if err != nil {
		return nil, 0, err;
	}

	if _, err := reader.ReadByte(); err != nil {
		return nil, 0, err;
	}
	return buf, size, err;
}

func EscapeMpd(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`);
	s = strings.ReplaceAll(s, `"`, `\"`);
	return `"` + s + `"`;
}

func GetCurrentSong() (Song, error) {
	resp, err := Request("currentsong");
	if err != nil {
		return Song{}, err;
	}
	songs, err := ParseInfoResponse(resp);
	if err != nil {
		return Song{}, err;
	}
	current := songs[0];

	return current, nil;
}

func ParseInfoResponse(plainResponse string) ([]Song, error) {
	var queue []Song;
	var s Song;

	for line := range strings.SplitSeq(plainResponse, "\n") {
		if strings.TrimSpace(line) == "" {
			continue;
		}
	
		parts := strings.SplitN(line, ": ", 2);
		if len(parts) != 2 {
			return nil, fmt.Errorf("failed to parse invalid line: `%v`", line);
		}
		key := parts[0];
		value := strings.TrimSpace(parts[1]);

		var conv int;
		var err error;

		switch key {
		case "file":
			if s.File != "" {
				queue = append(queue, s);
				s = Song{};
			}
			s.File = value;
		case "Artist":
			s.Artist = value;
		case "Title":
			s.Title = value;
		case "Album":
			s.Album = value;
		case "Pos":
			conv, err = strconv.Atoi(value);
			s.Pos = conv;
		case "Id":
			conv, err = strconv.Atoi(value);
			s.Id = conv;
		case "Time":
			conv, err = strconv.Atoi(value);
			s.Time = conv;
		case "Date":
			conv, err = strconv.Atoi(value);
			s.Date = conv;
		case "Duration":
			conv, err = strconv.Atoi(value);
			s.Duration = conv*60;
		}
		if err != nil {
			return nil, err;
		}
	}
	queue = append(queue, s);

	return queue, nil;
}


func ParseStatusResponse(plainResponse string) (Status, error) {
	var s Status;
	for line := range strings.SplitSeq(plainResponse, "\n") {
		if strings.TrimSpace(line) == "" {
			continue;
		}
		parts := strings.SplitN(line, ": ", 2);
		if len(parts) != 2 {
			return Status{}, fmt.Errorf("failed to parse invalid line: `%v`", line);
		}
		key := parts[0];
		value := strings.TrimSpace(parts[1]);

		var conv int;
		var convf float64;
		var err error;

		switch key {
		case "volume":
			conv, err = strconv.Atoi(value);
			s.Volume = conv;
		case "repeat":
			conv, err = strconv.Atoi(value);
			s.Repeat = conv;
		case "random":
			conv, err = strconv.Atoi(value);
			s.Random = conv;
		case "single":
			s.Single = value;
		case "consume":
			s.Consume = value;
		case "playlist":
			conv, err = strconv.Atoi(value);
			s.Playlist = conv;
	 	case "playlistlength":
			conv, err = strconv.Atoi(value);
			s.Playlistlength = conv;
	 	case "song":
			conv, err = strconv.Atoi(value);
			s.Song = conv;
	 	case "songid":
			conv, err = strconv.Atoi(value);
			s.SongId = conv;
	 	case "nextsong":
			conv, err = strconv.Atoi(value);
			s.NextSong = conv;
	 	case "nextsongid":
			conv, err = strconv.Atoi(value);
			s.NextSongId = conv;
		case "state":
			s.State = value;
		case "elapsed":
			convf, err = strconv.ParseFloat(value, 64);
			s.Elapsed = convf;
		case "duration":
			convf, err = strconv.ParseFloat(value, 64);
			s.Duration = convf;
		}
		if err != nil {
			return Status{}, err;
		}
	}
	return s, nil;
}

