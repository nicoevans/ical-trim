package parser

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var filter = Filter{"DTSTART", "greater_than", "20250426"}

type event struct {
	content []string
	fields  map[string]string
}

func Trim(iPath string) string {
	iFile, err := os.Open(iPath)
	if err != nil {
		log.Fatal("failed to open " + iPath)
	}
	defer iFile.Close()

	ext := filepath.Ext(iPath)
	oPath := strings.TrimSuffix(iPath, ext) + "_output" + ext
	oFile, err := os.Create(oPath)
	if err != nil {
		log.Fatal("failed to create " + oPath)
	}
	defer oFile.Close()

	reader := bufio.NewReader(iFile)
	writer := bufio.NewWriter(oFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			panic("unexpected error while reading ical: " + err.Error())
		}

		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			e := parseEvent(reader)
			if filter.shouldInclude(e.fields) {
				for _, l := range e.content {
					_, err = writer.Write([]byte(l + "\n"))
				}
			}
		} else {
			_, err = writer.Write([]byte(strings.TrimRight(line, " \n\r\t") + "\n"))
		}

		if err != nil {
			panic("unexpected error while reading ical: " + err.Error())
		}
	}

	writer.Flush()
	return oPath
}

func parseEvent(reader *bufio.Reader) event {
	content := []string{"BEGIN:VEVENT"}

	for {
		line, err := reader.ReadString('\n')

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
