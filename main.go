package main

import (
	"errors"
	"fmt"
	"mq/internal"
	"mq/internal/config"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	argv := os.Args[1:];
	err := parseCommandLineArguments(argv);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
		os.Exit(1); 
	}
}

func parseCommandLineArguments(argv []string) error {
	if len(argv) == 0 {
		return nil;
	}

	var err error;
	switch argv[0] {
	case "toggle":
		var command string = "pause";
		if len(argv) >= 2 {
			subcommand := argv[1];
			switch subcommand {
			case "consume", "single", "random", "repeat":
				command = subcommand;
			default:
				return fmt.Errorf("invalid subcommand to the `toggle` command: %v", subcommand);
			}
		}
		internal.ToggleCommand(command);
	case "stop":
		err = internal.RequestWithoutResponse("stop");
	case "prev":
		err = internal.RequestWithoutResponse("previous");
	case "next":
		err = internal.RequestWithoutResponse("next");
	case "delete":
		// TODO: support ranges
		if len(argv) < 2 {
			return errors.New("command `delete` needs a argument: song id");
		}
		req := fmt.Sprintf("delete %v", argv[1]);
		err = internal.RequestWithoutResponse(req);
	case "update":
		err = internal.RequestWithoutResponse("update");
	case "play":
		if len(argv) < 2 {
			return errors.New("command play needs a argument: song id");
		}
		req := fmt.Sprintf("play %v", argv[1]);
		err = internal.RequestWithoutResponse(req);
	case "add":
		if len(argv) < 2 {
			return errors.New("command add needs a argument: URI");
		}
		uri := argv[1];
		req := fmt.Sprintf("add %v", uri);
		err = internal.RequestWithoutResponse(req);
	case "see":
		config, err := config.GetConfig();
		if err != nil {
			return err;
		}
		basepath := config.BasePath;

		input := ".";
		if len(argv) == 2 {
			input = argv[1];
		}
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

	case "list", "ls":
		plainResp, err := internal.Request("playlistinfo");
		if err != nil {
			return err;
		}
		queue, err := internal.ParseInfoResponse(plainResp);
		if err != nil {
			return err;
		}
		if err := internal.PrintFormattedQueue(queue); err != nil {
			return err;
		}
	case "status":
		plainResp, err := internal.Request("status");
		if err != nil {
			return err;
		}
		if err := internal.PrintFormattedStatus(plainResp); err != nil {
			return err;
		}
	case "plain":
		if len(argv) < 2 {
			return errors.New("command add needs a argument: request");
		}
		request := argv[1];
		resp, err := internal.Request(request);
		if err != nil {
			return err;
		}
		fmt.Println(resp);


	default:
		return fmt.Errorf("command doesn't exist: %v", argv[0]);
	}

	if err != nil {
		return err;
	}
	
	return nil;
}


