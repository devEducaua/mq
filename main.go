package main

import (
	"os"
	"fmt"
	"errors"
	"mq/internal/mpd"
	"mq/internal/commands"
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
	var command string = "status";
	if len(argv) > 0 {
		command = argv[0];
	}

	var err error;
	switch command {
	case "toggle":
		togglecommand := "pause";
		if len(argv) >= 2 {
			subcommand := argv[1];
			switch subcommand {
			case "consume", "single", "random", "repeat":
				togglecommand = subcommand;
			default:
				return fmt.Errorf("invalid subcommand to the `toggle` command: %v", subcommand);
			}
		}
		err = commands.ToggleCommand(togglecommand);
	case "stop":
		err = mpd.RequestWithoutResponse("stop");
	case "clear":
		err = mpd.RequestWithoutResponse("clear");
	case "prev":
		err = mpd.RequestWithoutResponse("previous");
	case "next":
		err = mpd.RequestWithoutResponse("next");
	case "delete":
		// TODO: support ranges
		if len(argv) < 2 {
			return errors.New("command `delete` needs a argument: song id");
		}
		req := fmt.Sprintf("delete %v", argv[1]);
		err = mpd.RequestWithoutResponse(req);
	case "update":
		err = mpd.RequestWithoutResponse("update");
	case "play":
		if len(argv) < 2 {
			return errors.New("command play needs a argument: song id");
		}
		req := fmt.Sprintf("play %v", argv[1]);
		err = mpd.RequestWithoutResponse(req);
	case "add":
		if len(argv) < 2 {
			return errors.New("command add needs a argument: URI");
		}
		uri := argv[1];
		req := fmt.Sprintf("add %v", uri);
		err = mpd.RequestWithoutResponse(req);
	case "see":
		input := "/";
		if len(argv) == 2 {
			input = argv[1];
		}
		err = commands.SeeCommand(input);

	case "list", "ls":
		err = commands.ListCommand();
	case "status":
		plainResp, err := mpd.Request("status");
		if err != nil {
			return err;
		}
		if err := mpd.PrintFormattedStatus(plainResp); err != nil {
			return err;
		}
	case "plain":
		if len(argv) < 2 {
			return errors.New("command add needs a argument: request");
		}
		request := argv[1];
		resp, err := mpd.Request(request);
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


