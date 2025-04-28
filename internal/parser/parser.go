package parser

import (
	"bufio"
	"strings"
	"unicode"
)

var filter = Filter{"DTSTART", "greater_than", "20250426"}

type event struct {
	content []string
	fields  map[string]string
}

func Trim(r *bufio.Reader, w *bufio.Writer) {
	for {
		line, err := r.ReadString('\n')

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			panic("unexpected error while reading ical: " + err.Error())
		}

		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			e := parseEvent(r)
			if filter.shouldInclude(e.fields) {
				for _, l := range e.content {
					for len(l) > 1 {
						i := min(75, len(l))
						_, err = w.Write([]byte(l[:i] + "\r\n"))
						l = " " + l[i:]
					}
				}
			}
		} else {
			_, err = w.Write([]byte(line))
		}

		if err != nil {
			panic("unexpected error while reading ical: " + err.Error())
		}
	}

	w.Flush()
}

func parseEvent(r *bufio.Reader) event {
	content := []string{"BEGIN:VEVENT"}

	for {
		line, err := r.ReadString('\n')

		if err != nil {
			panic("unexpected error while reading event: " + err.Error())
		}

		if strings.HasPrefix(line, "END:VEVENT") {
			content = append(content, "END:VEVENT")
			break
		}

		if unicode.IsSpace(rune(line[0])) {
			content[len(content)-1] = content[len(content)-1] + strings.TrimSpace(line)
		} else {
			content = append(content, strings.TrimSpace(line))
		}
	}

	fields := make(map[string]string)
	for _, line := range content {
		pair := strings.SplitN(line, ":", 2)
		if len(pair) == 2 {
			fields[pair[0]] = pair[1]
		}
	}

	return event{content, fields}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
