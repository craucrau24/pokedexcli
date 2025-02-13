package main

import "strings"

func cleanInput(text string) []string {
	var result []string

	for _, item := range strings.Split(strings.Trim(text, " "), " ") {
		trimmed := strings.Trim(item, " ")
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
