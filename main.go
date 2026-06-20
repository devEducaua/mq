package main

import (
	"os"
	"fmt"
	"errors"
	"mq/internal/commands"
	"mq/internal/config"
	"mq/internal/mpd"
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
	config, err := config.GetConfig();
	if err != nil {
		return err;
	}

	var command string = config.DefaultCommand;
	if len(argv) > 0 {
		command = argv[0];
	}

	switch command {
	case "toggle":
		err = mpd.RequestWithoutResponse("pause");
	case "prev":
		err = mpd.RequestWithoutResponse("previous");
	case "stop", "clear", "next", "update":
		err = mpd.RequestWithoutResponse(command);
	case "consume", "single", "random", "repeat":
		err = commands.ChangeState(command);
	case "delete", "del":
		if len(argv) < 2 {
			return errors.New("command `delete` needs a argument: song id");
		}
		req := fmt.Sprintf("delete %v", argv[1]);
		err = mpd.RequestWithoutResponse(req);
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
		uri := mpd.EscapeMpd(argv[1]);
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
	case "move", "mv":
		if len(argv) < 3 {
			return fmt.Errorf("command `%v` needs two arguments: `from` and `to` ", command);
		}
		from := argv[1];
		to := argv[2];
		req := fmt.Sprintf("move %v %v", from, to);
		err = mpd.RequestWithoutResponse(req);
	case "status":
		plainResp, err := mpd.Request("status");
		if err != nil {
			return err;
		}
		if err := mpd.PrintFormattedStatus(plainResp); err != nil {
			return err;
		}
	case "albumart":
		notify := false;	
		if len(argv) <= 2 {
			if argv[1] == "--notify" {
				notify = true;
			}
		}
		err = commands.AlbumArt(notify);
	case "--help":
		err = commands.RunExternalCommand(true, "man", "mq(1)");
	default:
		return fmt.Errorf("command doesn't exist: %v", argv[0]);
	}

	if err != nil {
		return err;
	}
	
	return nil;
}


