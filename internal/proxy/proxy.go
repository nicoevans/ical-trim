package proxy

import (
	"net/http"
	"os"

	"github.com/nicoevans/ical-trim/internal/parser"
)

func Proxy(url string) {
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	parser.Trim(resp.Body, os.Stdout)
}
