package commands

import (
	"bufio"
	"fmt"
	"mq/internal/config"
	"mq/internal/mpd"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ChangeState(state string, newMode string) error {
	var mode string;
	mode = newMode;

	if newMode == "" {
		plainResponse, err := mpd.Request("status");
		if err != nil {
			return err;
		}
		status, err := mpd.ParseStatusResponse(plainResponse);
		if err != nil {
			return err;
		}

		switch state {
		case "repeat":
			mode = fromIntToState(1 - status.Repeat);
		case "random":
			mode = fromIntToState(1 - status.Random);
		case "consume":
			mode = fromIntToState(1 - int(status.Consume));
		case "single":
			mode = fromIntToState(1 - int(status.Single));
		}
	}

	request := fmt.Sprintf("%v %v", state, mode);
	if err := mpd.RequestWithoutResponse(request); err != nil {
		return err;
	}
	return nil;
}

func fromIntToState(i int) string {
	switch i {
	case 0, 1:
		return strconv.Itoa(i);
	case -1:
		return "oneshot";
	default:
		return "";
	}
}

func SeeCommand(input string) error {
	config, err := config.GetConfig();
	if err != nil {
		return err;
	}
	basepath := config.BasePath;

	path := filepath.Join(basepath, input);
	entries, err := os.ReadDir(path);
	if err != nil {
		return err;
	}
	for _,e := range entries {
		name := e.Name();
		if strings.HasPrefix(name, ".") {
			continue;
		}
		fmt.Printf("%v\n", name);
	}
	return nil;
}

func ListCommand() error {
	plainResp, err := mpd.Request("playlistinfo");
	if err != nil {
		return err;
	}
	queue, err := mpd.ParseInfoResponse(plainResp);
	if err != nil {
		return err;
	}
	if err := mpd.PrintFormattedQueue(queue); err != nil {
		return err;
	}
	return nil;
}

func AlbumArt(notify bool) error {
	msgs := make(chan string);
	errs := make(chan error);

	go watchPlayer(msgs, errs);

	for {
		select {
		case msg := <-msgs:
			if msg == "changed: player\n" {
				if err := writeImageToPath(); err != nil {
					return err;
				}
				if notify {
					runNotify();
				}
			}
		case err := <-errs:
			return err;
		}
	}
}

func runNotify() error {
	s, err := mpd.GetCurrentSong();
	if err != nil {
		return err;
	}

	config, err := config.GetConfig()
	if err != nil {
		return err;
	}

	if err := RunExternalCommand(false, config.NotifyScriptPath, config.CoverOutputPath, s.Artist, s.Title); err != nil {
		return err;
	}

	return nil;
}

func RunExternalCommand(setStds bool, command ...string) error {
	cmd := exec.Command(command[0], command[1:]...);

	if setStds {
		cmd.Stdin = os.Stdin;
		cmd.Stdout = os.Stdout;
		cmd.Stderr = os.Stderr;
	}
	if err := cmd.Run(); err != nil {
		return err;
	}

	return nil;
}

func writeImageToPath() error {
	config, err := config.GetConfig();
	if err != nil {
		return err;
	}

	outputPath := config.CoverOutputPath;

	s, err := mpd.GetCurrentSong();
	if err != nil {
		return err;
	}
	if s.File == "" {
		return nil;
	}

	image, found, err := findCoverFile(s.File);
	if err != nil {
		return err;
	}
	
	if !found {
		image, err = readPictureFile(s.File);
		if err != nil {
			return err;
		}
	}
	if err := os.WriteFile(outputPath, image, 0755); err != nil {
		return err;
	}

	return nil;
}

func findCoverFile(currentFile string) ([]byte, bool, error) {
	config, err := config.GetConfig();
	if err != nil {
		return nil, false, err;
	}

	file := filepath.Join(config.BasePath, currentFile);

	found := false;
	var image []byte;

	dir := filepath.Dir(file);
	entries, err := os.ReadDir(dir);
	if err != nil {
		return nil, false, fmt.Errorf("cannot open directory: %v", dir);
	}

	for _,e := range entries {
		name := e.Name();
		if name == "cover.jpg" || name == "cover.png" || name == "cover.jxl" || name == "cover.webp" {
			found = true;
			path := filepath.Join(dir, name);
			image, err = os.ReadFile(path);
			if err != nil {
				return nil, false, err;
			}
		}
	}
	return image, found, nil;
}

func readPictureFile(file string) ([]byte, error) {
	escaped := mpd.EscapeMpd(file);
	var image []byte;
	var size, offset int;

	for {
		req := fmt.Sprintf("readpicture %v %v", escaped, offset);
		resp, err := mpd.Request(req);
		if err != nil {
			return nil, err;
		}

		reader := bufio.NewReader(strings.NewReader(resp));

		chunk, totalSize, err := mpd.ParseBinaryResponse(reader);
		if err != nil {
			return nil, err;
		}

		size = totalSize;
		image = append(image, chunk...);
		offset += len(chunk);

		if len(image) >= size {
			break;
		}
	}
	return image, nil;
}

func watchPlayer(msg chan<- string, errs chan<- error) {
	for {
		result, err := mpd.Request("idle player");
		if err != nil {
			errs <- err;
			continue;
		}
		msg <- result;
	}
}

