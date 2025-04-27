package main

import (
	"os"

	"github.com/nicoevans/ical-trim/internal/parser"
)

func main() {
	parser.Trim(os.Args[1])
}
