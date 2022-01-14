package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

type Config struct {
	Level    string
	Server   map[string]ServerConfig
	DebugSer *string
}

type ServerConfig struct {
	Local  string
	Remote string
	Type   string
}

func getConfigPath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, execname := filepath.Split(fullexecpath)
	ext := filepath.Ext(execname)
	name := execname[:len(execname)-len(ext)]
	fname := filepath.Join(dir, name+".yaml")

	files, err := filepath.Glob(fname)
	if err == nil && len(files) > 0 {
		return files[0], nil
	}
	files, err = filepath.Glob(filepath.Join(dir, name+".json"))
	if err == nil && len(files) > 0 {
		return files[0], nil
	}
	return fname, nil
}

func GetConfig() (*Config, error) {
	f, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	cfg := &Config{Server: make(map[string]ServerConfig)}
	err = configor.Load(cfg, f)
	if err != nil {
		log.Println(err)
	}
	return cfg, nil
}
