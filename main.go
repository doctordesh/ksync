package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	LocalPath  string   `json:"local_path"`
	RemotePath string   `json:"remote_path"`
	User       string   `json:"user"`
	Host       string   `json:"host"`
	Exclude    []string `json:"exclude"`
}

func getConfig() (Config, error) {
	var config Config
	var err error

	urs, err := user.Current()
	if err != nil {
		return config, err
	}

	configFilePath := filepath.Join(urs.HomeDir, ".ksync")
	f, err := os.Open(configFilePath)
	if err != nil {
		return config, fmt.Errorf("ERROR: Config file ~/.ksync is missing")
	}

	configFile, err := ioutil.ReadAll(f)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return config, fmt.Errorf("ERROR: ~/.ksync is not well formatted JSON (%s)", err)
	}

	if config.LocalPath == "" {
		return config, fmt.Errorf("ERROR: Key 'local_path' missing from ~/.ksync config file")
	}

	if config.RemotePath == "" {
		return config, fmt.Errorf("ERROR: Key 'remote_path' missing from ~/.ksync config file")
	}

	if config.User == "" {
		return config, fmt.Errorf("ERROR: Key 'user' missing from ~/.ksync config file")
	}

	if config.Host == "" {
		return config, fmt.Errorf("ERROR: Key 'host' missing from ~/.ksync config file")
	}

	return config, nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	relativePath, err := filepath.Rel(config.LocalPath, dir)
	if err != nil {
		panic(err)
	}

	if strings.Contains(relativePath, "..") {
		fmt.Println("ERROR: You're not in the base path")
		return
	}

	rsyncFromPath := path.Join(config.LocalPath, relativePath)
	rsyncToPath := path.Join(config.RemotePath, relativePath)

	var args []string

	for _, exclude := range config.Exclude {
		args = append(args, fmt.Sprintf("--exclude=\"%s\"", exclude))
	}

	args = append(
		args,
		"-azP",
		"--delete",
		rsyncFromPath,
		fmt.Sprintf("%s@%s:%s", config.User, config.Host, rsyncToPath),
	)

	cmd := exec.Command("rsync", args...)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
