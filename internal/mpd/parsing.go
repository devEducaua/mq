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

type State int;
const (
	StateOff State = iota
	StateOn
	StateOneshot
)

type Status struct {
	Volume int
	Repeat int
	Random int
	Single State
	Consume State
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
			conv, err := strconv.Atoi(value);
			if err != nil {
				return nil, err;
			}
			s.Pos = conv;
		case "Id":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return nil, err;
			}
			s.Id = conv;
		case "Time":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return nil, err;
			}
			s.Time = conv;
		case "Date":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return nil, err;
			}
			s.Date = conv;
		case "Duration":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return nil, err;
			}
			s.Duration = conv*60;
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

		switch key {
		case "volume":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.Volume = conv;
		case "repeat":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.Repeat = conv;
		case "random":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.Random = conv;
		case "single":
			var v State;
			switch value {
			case "0":
				v = StateOff;
			case "1":
				v = StateOn;
			case "oneshot":
				v = StateOneshot;
			default:
				return Status{}, fmt.Errorf("invalid single state");
			}
			s.Single = v;
		case "consume":
			var v State;
			switch value {
			case "0":
				v = StateOff;
			case "1":
				v = StateOn;
			case "oneshot":
				v = StateOneshot;
			default:
				return Status{}, fmt.Errorf("invalid consume state");
			}
			s.Consume = v;
		case "playlist":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.Playlist = conv;
	 	case "playlistlength":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.Playlistlength = conv;
	 	case "song":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.Song = conv;

	 	case "songid":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.SongId = conv;

	 	case "nextsong":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.NextSong = conv;

	 	case "nextsongid":
			conv, err := strconv.Atoi(value);
			if err != nil {
				return Status{}, err;
			}
			s.NextSongId = conv;
		case "state":
			s.State = value;
		case "elapsed":
			conv, err := strconv.ParseFloat(value, 64);
			if err != nil {
				return Status{}, err;
			}
			s.Elapsed = conv;
		case "duration":
			conv, err := strconv.ParseFloat(value, 64);
			if err != nil {
				return Status{}, err;
			}
			s.Duration = conv;
		}
	}
	return s, nil;
}

