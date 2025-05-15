package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gitlab.maxiv.lu.se/emiros/ksync"
)

func usage() {
	fmt.Println("Usage: ksync [options]")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("KSYNC_CONFIG environment variable can be used to select config file. Default value is ~/.ksync")
}

func main() {
	flag.Parse()

	configPath := os.Getenv("KSYNC_CONFIG")
	if configPath == "" {

		urs, err := user.Current()
		if err != nil {
			panic(err)
		}

		configPath = filepath.Join(urs.HomeDir, ".ksync")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = ksync.Ksync(cwd, configPath)
	if err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
