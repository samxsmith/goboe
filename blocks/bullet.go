package blocks

import (
	"regexp"
	"strings"
)

var bulletWithoutNewLine = regexp.MustCompile("\n-")

// FormatBullets ensures bullets will appear correctly.
// The parser expects a new line before bullets.
func FormatBullets(noteBody string) string {
	//var lines []string

	lines := strings.Split(noteBody, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if !strings.HasPrefix(line, "-") {
			continue
		}

		previousLine := lines[i-1]
		if strings.HasPrefix(previousLine, "-") {
			// if previous line was bullet, skip
			continue
		}

		if previousLine == "" {
			// if previous line was empty, that's perfect
			continue
		}

		// add newline after previous line
		lines[i-1] = lines[i-1] + "\n"
	}

	return strings.Join(lines, "\n")
}
