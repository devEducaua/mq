package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	argv := os.Args[1:];
	if len(argv) == 0 {
		os.Exit(1);
	}

	if argv[0] == "toggle" {
		err := request("pause");
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
			os.Exit(1);
		}
	}
}

func initializeMpdConnection() (net.Conn, error) {
	var conn net.Conn;

	addr, err := net.ResolveTCPAddr("tcp", "localhost:6600");
	if err != nil {
		return conn, err;
	}

	conn, err = net.DialTCP("tcp", nil, addr);
	if err != nil {
		return conn, err;
	}

	reader := bufio.NewReader(conn);
	line, err := reader.ReadString('\n');
	if err != nil {
		return conn, err;
	}

	if !strings.HasPrefix(line, "OK MPD") {
		return conn, fmt.Errorf("failed to initialize mpd connection");
	}

	return conn, nil;
}

func request(request string) (error) {
	conn, err := initializeMpdConnection();
	if err != nil {
		return err;
	}
	if conn != nil {
		defer conn.Close();
	}

	var line string;
	var reader = bufio.NewReader(conn);
	var sb strings.Builder;

	fmt.Fprintf(conn, "%v\n", request);
	for {
		if line, err = reader.ReadString('\n'); err != nil {
			return err;
		}
		if line == "OK\n" {
			break;
		} else if strings.HasPrefix(line, "ACK ") {
			err = fmt.Errorf("ERROR: request failed: %v", sb.String());
			break;
		} else {
			sb.WriteString(line);
		}
	}

	return err;
}
