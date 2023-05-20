// Package models provides a collection of struct types for unmarshaling command line flags.
package models

//
// This package includes the Flags struct, which represents the command line flags that can be passed as arguments.
// It allows users to easily access and manipulate the values of the command line flags.
//

// Flags represents the command line flags that are unmarshaled into struct fields.
type Flags struct {
	Word   string
	Input  string
	Output string
	Level  int
}
