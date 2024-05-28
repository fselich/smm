package ui

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/muesli/reflow/wordwrap"
	"regexp"
	"strings"
)

func SyntaxHighlight(secretData []byte) string {
	var buf bytes.Buffer
	format := detectFormat(secretData)
	err := quick.Highlight(&buf, string(secretData), format, "terminal", "rrt")
	if err != nil {
		panic(err)
	}

	text := wordwrap.String(buf.String(), 255)
	return text
}

func detectFormat(secretData []byte) string {
	if isEnvFormat(string(secretData)) {
		return "bash"
	}
	if isJSON(string(secretData)) {
		return "json"
	}

	return "default"
}

func isEnvFormat(s string) bool {
	lines := strings.Split(s, "\n")

	envRegex := regexp.MustCompile(`^\s*[a-zA-Z_][a-zA-Z0-9_]*\s*=\s*.*$`)

	for _, line := range lines {
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		if !envRegex.MatchString(line) {
			return false
		}
	}

	return true
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
