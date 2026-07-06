package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	config := configFromArgs()

	err := process(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unfucking: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Done\n")

}
