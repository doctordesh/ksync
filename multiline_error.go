package ksync

import "strings"

type MultilineError struct {
	Parts []error
}

// Error ...
func (self MultilineError) Error() string {
	parts := []string{}
	for _, err := range self.Parts {
		parts = append(parts, err.Error())
	}

	return strings.Join(parts, "\n")
}

// Add ...
func (self *MultilineError) Add(err error) {
	self.Parts = append(self.Parts, err)
}
