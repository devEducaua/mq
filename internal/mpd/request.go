package mpd

import (
	"bufio"
	"fmt"
	"io"
	"mq/internal/config"
	"net"
	"strings"
)

func initializeMpdConnection() (*net.TCPConn, error) {
	config, err := config.GetConfig();
	if err != nil {
		return nil, err;
	}

	addr, err := net.ResolveTCPAddr("tcp", config.Addr);
	if err != nil {
		return nil, err;
	}

	conn, err := net.DialTCP("tcp", nil, addr);
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mpd: %v", err);
	}

	reader := bufio.NewReader(conn);
	line, err := reader.ReadString('\n');
	if err != nil {
		return nil, err;
	}

	if !strings.HasPrefix(line, "OK MPD") {
		return nil, fmt.Errorf("failed to initialize mpd connection");
	}

	return conn, nil;
}

func Request(request string) (string, error) {
	conn, err := initializeMpdConnection();
	if err != nil {
		return "", err;
	}
	if conn != nil {
		defer conn.Close();
	}

	var reader = bufio.NewReader(conn);
	var sb strings.Builder;

	fmt.Fprintf(conn, "%v\n", request);
	for {
		line, err := reader.ReadString('\n'); 
		if err != nil {
			if err == io.EOF {
				return "",  io.ErrUnexpectedEOF;
			}
			return "",  fmt.Errorf("failed to read line: %v", err);
		}
		if line == "OK\n" {
			break;
		} else if strings.HasPrefix(line, "ACK ") {
			return "", fmt.Errorf("request failed: %v", line);
		} else {
			sb.WriteString(line);
		}
	}

	return sb.String(), nil;
}

func RequestWithoutResponse(request string) (error) {
	_, err := Request(request);
	if err != nil {
		return err;
	}
	return nil;
}

