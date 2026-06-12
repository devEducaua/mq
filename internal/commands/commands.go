package commands

import (
	"bufio"
	"fmt"
	"mq/internal/config"
	"mq/internal/mpd"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ToggleCommand(command string) error {
	if command != "pause" {
		plainResponse, err := mpd.Request("status");
		if err != nil {
			return err;
		}
		status, err := mpd.ParseStatusResponse(plainResponse);
		if err != nil {
			return err;
		}
		var mode int;

		switch command {
		case "repeat":
			mode = 1 - status.Repeat;
		case "random":
			mode = 1 - status.Random;
		case "single":
			mode = 1 - status.Single;
		case "consume":
			mode = 1 - status.Consume;
		default:
			return fmt.Errorf("invalid subcommand to the `toggle` command: %v", command);
		}
		command = fmt.Sprintf("%v %v", command, mode);
	}

	if err := mpd.RequestWithoutResponse(command); err != nil {
		return err;
	}
	return nil;
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
		if !strings.HasPrefix(name, ".") {
			fmt.Println(name);
		}
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

func AlbumArt() error {

	msgs := make(chan string);
	errs := make(chan error);

	go watchPlayer(msgs, errs);

	for {
		select {
		case msg := <-msgs:
			if msg == "changed: player\n" {
				path, err := getAlbumPath();
				if err != nil {
					return err;
				}
				runImageCommand(path);
			}
		case err := <-errs:
			return err;
		}
	}
}

var imageCmd *exec.Cmd;
func runImageCommand(path string) error {

	if imageCmd != nil && imageCmd.Process != nil {
		if err := imageCmd.Process.Kill(); err != nil {
			return err;
		}
		_, err := imageCmd.Process.Wait();
		if err != nil {
			return err;
		}
	}

	config, err := config.GetConfig();
	if err != nil {
		return err;
	}
	parts := strings.Fields(config.ImageCommand);
	imageCmd = exec.Command(parts[0], append(parts[1:], path)...)

	if err := imageCmd.Start(); err != nil {
		return err;
	}

	return nil;
}

func getAlbumPath() (string, error) {
	s, err := mpd.GetCurrentSong();
	if err != nil {
		return "", err;
	}

	config, err := config.GetConfig();
	if err != nil {
		return "", err;
	}

	file := filepath.Join(config.BasePath, s.File);

	dir := filepath.Dir(file);
	entries, err := os.ReadDir(dir);
	if err != nil {
		return "", fmt.Errorf("cannot open directory: %v", dir);
	}
	for _,e := range entries {
		name := e.Name();
		if name == "cover.jpg" || name == "cover.png" || name == "cover.jxl" || name == "cover.webp" {
			return filepath.Join(dir, name), nil;
		}
	}

	escaped := mpd.EscapeMpd(s.File);
	var image []byte;
	var size, offset int;

	for {
		req := fmt.Sprintf("readpicture %v %v", escaped, offset);
		resp, err := mpd.Request(req);
		if err != nil {
			return "", err;
		}

		reader := bufio.NewReader(strings.NewReader(resp));

		chunk, totalSize, err := mpd.ParseBinaryResponse(reader);
		if err != nil {
			return "", err;
		}

		size = totalSize;
		image = append(image, chunk...);
		offset += len(chunk);

		if len(image) >= size {
			break;
		}
	}

	path := "/tmp/cover.jpg";
	if err := os.WriteFile(path, image, 0755); err != nil {
		return "", err;
	}

	return path, nil;
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
