package ksync

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func runSync(cwd string, config Config, target Target, wg *sync.WaitGroup) {
	defer func(wg *sync.WaitGroup) {
		wg.Done()
	}(wg)

	var err error

	// This was done already, just need the value
	rel, _ := filepath.Rel(target.Source, cwd)

	rsyncFromPath := path.Clean(cwd) + "/"
	rsyncToPath := path.Join(target.Destination.Path, rel)

	if config.Verbose {
		fmt.Printf("using target '%s'\n", target.Name)
		fmt.Printf("   from path %s\n", rsyncFromPath)
		fmt.Printf("     to path %s@%s:%s\n", target.Destination.User, target.Destination.Host, rsyncToPath)
		fmt.Println()
	}

	// ==================================================
	// Create all directories

	// args := []string{
	// 	fmt.Sprintf("%s@%s", target.Destination.User, target.Destination.Host),
	// 	fmt.Sprintf("'mkdir -p %s'", rsyncToPath),
	// }

	// cmd := exec.Command("ssh", args...)
	// fmt.Printf("COMMAND: %+v\n\n\n", cmd)
	// if config.Verbose {
	// 	fmt.Printf("running: ssh %s\n", strings.Join(args, " "))
	// 	fmt.Println("########## SSH OUTPUT ##########")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// }

	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Printf("ERROR: %s\n", err)
	// 	return
	// }

	// if config.Verbose {
	// 	fmt.Println("##########     END      ##########")
	// 	fmt.Println()
	// }

	// ==================================================
	// RSync content

	args := []string{
		"-azP",
		"--delete",
		rsyncFromPath,
		fmt.Sprintf("%s@%s:%s", target.Destination.User, target.Destination.Host, rsyncToPath),
	}

	cmd := exec.Command("rsync", args...)
	cmd.Stderr = os.Stderr

	if config.Verbose {
		fmt.Printf("running: rsync %s\n", strings.Join(args, " "))
		fmt.Println("########## RSYNC OUTPUT ##########")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	if config.Verbose {
		fmt.Println("##########     END      ##########")
		fmt.Println()
	}
}

// ==================================================
//
// Shell
//
// Input: flags and environment variables
//
// Output: <nothing>
//
// ==================================================

func Ksync(cwd, configPath string) error {
	var err error
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	config, err := loadConfig(f)
	if err != nil {
		return fmt.Errorf("file %s; %w", configPath, err)
	}

	errs := config.Validate()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
		return fmt.Errorf("invalid config file at %s", configPath)
	}

	targets, err := findTargets(config.Targets, cwd)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)
		go runSync(cwd, config, target, &wg)
		// // This was done already, just need the value
		// rel, _ := filepath.Rel(target.Source, cwd)

		// rsyncFromPath := path.Clean(cwd) + "/"
		// rsyncToPath := path.Join(target.Destination.Path, rel)

		// if config.Verbose {
		// 	fmt.Printf("using target '%s'\n", target.Name)
		// 	fmt.Printf("   from path %s\n", rsyncFromPath)
		// 	fmt.Printf("     to path %s@%s:%s\n", target.Destination.User, target.Destination.Host, rsyncToPath)
		// 	fmt.Println()
		// }

		// args := []string{
		// 	"-azP",
		// 	"--delete",
		// 	rsyncFromPath,
		// 	fmt.Sprintf("%s@%s:%s", target.Destination.User, target.Destination.Host, rsyncToPath),
		// }

		// cmd := exec.Command("rsync", args...)
		// cmd.Stderr = os.Stderr

		// if config.Verbose {
		// 	fmt.Printf("running: rsync %s\n", strings.Join(args, " "))
		// 	fmt.Println("########## RSYNC OUTPUT ##########")
		// 	cmd.Stdout = os.Stdout
		// }

		// err = cmd.Run()
		// if err != nil {
		// 	return err
		// }

		// if config.Verbose {
		// 	fmt.Println("##########     END      ##########")
		// 	fmt.Println()
		// }
	}

	wg.Wait()

	return nil
}

// ==================================================
//
// Core
//
// Input: path to sync, config
//
// Output: args for a rsync command
//
// ==================================================

func findTargets(targets Targets, cwd string) ([]Target, error) {
	var found []Target

	for _, tar := range targets {
		rel, err := filepath.Rel(tar.Source, cwd)
		if err != nil {
			panic(err)
		}

		// If the rel want's to go /out/ of the directory, then we're
		// not in te same base path.
		if strings.Contains(rel, "..") == false {
			found = append(found, tar)
		}
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("path %s is not within any target", cwd)
	}

	return found, nil
}

func loadConfig(stream io.Reader) (Config, error) {
	var err error
	var config Config

	err = json.NewDecoder(stream).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("could not parse config file: %w", err)
	}

	return config, nil
}
