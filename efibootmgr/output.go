package efibootmgr

import (
	"fmt"
	"regexp"
	"strings"
)

type Output struct {
	BootCurrent string
	Timeout     string
	BootOrder   []string
	BootEntries []BootEntry
}

type BootEntry struct {
	Bootnum string
	Label   string
}

func ParseOutput(output string) *Output {
	o := &Output{}

	bootCurrentRegex := regexp.MustCompile("^BootCurrent: (\\d{4})")
	timeoutRegex := regexp.MustCompile("^Timeout: (.*)")
	bootOrderRegex := regexp.MustCompile("^BootOrder: ([0-9,]*)")
	noBootOrderSetRegex := regexp.MustCompile("^No BootOrder is set;")
	bootEntryRegex := regexp.MustCompile("^Boot(\\d{4})([* ]) ([^\\t]*)\\t")

	for _, line := range strings.Split(output, "\n") {
		if bootCurrentRegex.MatchString(line) {
			m := bootCurrentRegex.FindAllStringSubmatch(line, -1)
			o.BootCurrent = m[0][1]

			continue
		}
		if timeoutRegex.MatchString(line) {
			m := timeoutRegex.FindAllStringSubmatch(line, -1)
			o.Timeout = m[0][1]

			continue
		}
		if bootOrderRegex.MatchString(line) {
			m := bootOrderRegex.FindAllStringSubmatch(line, -1)
			o.BootOrder = strings.Split(m[0][1], ",")

			continue
		}
		if noBootOrderSetRegex.MatchString(line) {
			continue
		}
		if bootEntryRegex.MatchString(line) {
			m := bootEntryRegex.FindAllStringSubmatch(line, -1)
			entry := BootEntry{
				Bootnum: m[0][1],
				Label:   m[0][3],
			}
			o.BootEntries = append(o.BootEntries, entry)

			continue
		}

		fmt.Printf("Could not parse line '%s'\n", line)
	}

	return o
}
