package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Type checks an interface to see if it is a string.
func IsString(candidate interface{}) bool {
	switch c := candidate.(type) {
	case string:
		var _ = c
		return true
	}

	return false
}

// Checks if a string to interface map contains a key.
func HasKey(key string, m map[string]interface{}) bool {
	_, ok := m[key]
	return ok
}

// Return a random integer within a range.
func RandRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// Print usage with an error code.
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(127)
}
