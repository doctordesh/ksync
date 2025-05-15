package ksync

import "fmt"

type Config struct {
	Verbose bool    `json:"verbose"`
	Targets Targets `json:"targets"`
}

type Targets []Target

type Target struct {
	Name        string      `json:"name"`
	Source      string      `json:"source"`
	Destination Destination `json:"destination"`
}

type Destination struct {
	User string `json:"user"`
	Host string `json:"host"`
	Path string `json:"path"`
}

// Validate ...
func (self Config) Validate() []error {
	var errs []error

	for i, c := range self.Targets {
		if c.Name == "" {
			errs = append(errs, fmt.Errorf("missing key 'name' for target %d", i))
			continue
		}

		if c.Source == "" {
			errs = append(errs, fmt.Errorf("target '%s' is missing key 'source'", c.Name))
		}

		if c.Destination.User == "" {
			errs = append(errs, fmt.Errorf("target '%s' is missing key 'destination.user'", c.Name))
		}

		if c.Destination.Host == "" {
			errs = append(errs, fmt.Errorf("target '%s' is missing key 'destination.host'", c.Name))
		}

		if c.Destination.Path == "" {
			errs = append(errs, fmt.Errorf("target '%s' is missing key 'destination.path'", c.Name))
		}
	}

	return errs
}
