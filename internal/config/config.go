package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Addr string
	BasePath string
}

func GetConfig() (Config, error) {
	var c Config;

	base, err := GetBaseDir();
	if err != nil {
		return Config{}, err;
	}

	path := filepath.Join(base, "config");

	dat, err := os.ReadFile(path);
	if err != nil {
		return Config{}, err;
	}
	content := string(dat);

	lines := strings.Split(content, "\n");
	for i := range lines {
		line := strings.TrimSpace(lines[i]);
		if line == "" || strings.HasPrefix(line, "#") {
			continue;
		}
		parts := strings.SplitN(line, ":", 2);
		key := strings.TrimSpace(parts[0]);
		value := strings.TrimSpace(parts[1]);

		switch key {
		case "addr":
			c.Addr = value;
		case "basepath":
			home, err := os.UserHomeDir();
			if err != nil {
				return Config{}, err;
			}
			var basepath string = value;
			if strings.HasPrefix(value, "~/") {
				basepath = filepath.Join(home, basepath[2:]);
			}
			c.BasePath = basepath;
		default:
			return Config{}, fmt.Errorf("fail to parse line in the config: %v", i);
		}
	}

	return c, nil;
}

func GetBaseDir() (string, error) {

	xdgDir := os.Getenv("XDG_CONFIG_HOME");

	path := filepath.Join(xdgDir, "mq");
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err;
	}

	return path, nil;
}

