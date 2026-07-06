package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	InPlace bool
	InPath  string
	OutPath string
}

func configFromArgs() Config {
	var inPlace bool
	var showVersion bool

	flag.BoolVar(&inPlace, "i", false, "Modify the input file in-place (temp-and-swap)")
	flag.BoolVar(&showVersion, "version", false, "Print version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_file> <output_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s -i <input_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("unfuck-properties v%s\n", version)
		os.Exit(0)
	}

	args := flag.Args()

	if inPlace {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Error: -i requires exactly 1 argument (the input file).")
			flag.Usage()
			os.Exit(1)
		}
		return Config{
			InPlace: true,
			InPath:  args[0],
		}
	}
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Requires exactly 2 arguments (input_file and output_file) when not using -i.")
		flag.Usage()
		os.Exit(1)
	}

	return Config{
		InPlace: false,
		InPath:  args[0],
		OutPath: args[1],
	}
}
