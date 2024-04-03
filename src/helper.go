package main

import (
	"os"
	"strings"
)

func GetToken(filename string) (token string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(wd + `\\` + filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func sanitizeInput(input string) string {
	// Replace single quotes with two single quotes to escape them
	sanitized := strings.Replace(input, "'", "''", -1)
	return sanitized
}
