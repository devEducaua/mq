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
	DefaultCommand string
	CoverOutputPath string
}

func GetConfig() (Config, error) {
	config := Config{
		Addr: "localhost:6600",
		BasePath: "./",
		DefaultCommand: "list",
		CoverOutputPath: "/tmp/cover.jpg",
	}

	parsed, err := ParseConfig();
	if err != nil {
		return Config{}, err;
	}

	if config.Addr != parsed.Addr && parsed.Addr != "" {
		config.Addr = parsed.Addr;	
	}

	if config.BasePath != parsed.BasePath && parsed.BasePath != "" {
		config.BasePath = parsed.BasePath;	
	}

	if config.DefaultCommand != parsed.DefaultCommand && parsed.DefaultCommand != "" {
		config.DefaultCommand = parsed.DefaultCommand;	
	}

	if config.CoverOutputPath != parsed.CoverOutputPath && parsed.CoverOutputPath != "" {
		config.CoverOutputPath = parsed.CoverOutputPath;	
	}

	return config, nil;
}

func ParseConfig() (Config, error) {
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
		case "base-path":
			home, err := os.UserHomeDir();
			if err != nil {
				return Config{}, err;
			}
			var basepath string = value;
			if strings.HasPrefix(value, "~/") {
				basepath = filepath.Join(home, basepath[2:]);
			}
			c.BasePath = basepath;
		case "default-command":
			c.DefaultCommand = value;
		case "cover-path":
			c.CoverOutputPath = value;
		default:
			return Config{}, fmt.Errorf("fail to parse line in the config: %v\n\t%v\n", i+1, line);
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

