package main

import (
	"flag"
)

func ParseInput() (string, string) {
	var (
		input  = flag.String("input", "", "Path to `.toml` input file.")
		output = flag.String("output", "", "Patho to `.tex` output file")
	)
	flag.Parse()

	return *input, *output
}
