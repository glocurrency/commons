package logger

import (
	"os"
	"regexp"
	"slices"
)

// Writer writes to stdout or stderr depending on severity level.
type Writer struct{}

func (Writer) Write(p []byte) (n int, err error) {
	if IsErrLog(p) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

var levelRegexT = regexp.MustCompile("level=([a-z]+)")
var levelRegexJ = regexp.MustCompile(`\"level\":\"([a-z]+)\"`)

func IsErrLog(p []byte) bool {
	must := []string{"error", "warning", "fatal", "panic"}

	matches := levelRegexJ.FindStringSubmatch(string(p))
	if len(matches) > 1 {
		return slices.Contains(must, matches[1])
	}

	matches = levelRegexT.FindStringSubmatch(string(p))
	if len(matches) > 1 {
		return slices.Contains(must, matches[1])
	}

	return false
}
