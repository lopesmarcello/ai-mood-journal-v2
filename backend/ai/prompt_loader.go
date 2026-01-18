package ai

import (
	"os"
	"strings"
)

func LoadSystemPrompt(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(data)), nil
}
