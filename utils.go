package main

import (
	"os/exec"
	"regexp"
	"strings"
)

func FormatName(name string) string {
	// Transforme camelCase ou snake_case en capitalisé
	// Exemple : helloWorld -> Hello World
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	name = re.ReplaceAllString(name, "$1 $2")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.Title(name)
	return name
}

func GetGitHash() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

