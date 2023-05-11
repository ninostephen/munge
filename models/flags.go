// Package models hold the struct for unmarshelling command line flags
package models

// Flags contains various fields to hold data passed as arguments
type Flags struct {
	Word   string
	Input  string
	Output string
	Level  int
}
