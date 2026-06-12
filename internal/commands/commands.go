package commands

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"mq/internal/config"
	"mq/internal/mpd"
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

