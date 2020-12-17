package main

import (
	"errors"
	"strings"
)

func FindInArrayLax(haystack []string, needle string) (string, error) {
	for _, element := range haystack {
		if strings.TrimSpace(strings.ToLower(element)) == needle {
			return needle, nil
		}
	}
	return "", errors.New("not found")
}
