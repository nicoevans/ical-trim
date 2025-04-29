package main

import (
	"github.com/nicoevans/ical-trim/internal/proxy"
)

func main() {
	proxy.Proxy("https://calendar.google.com/calendar/ical/restofsecreturl")
}
