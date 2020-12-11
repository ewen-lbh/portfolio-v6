package main

import "errors"

func FindInArray(haystack []string, needle string) (string, error) {
	for _, element := range haystack {
		if element == needle {
			return needle, nil
		}
	}
	return "", errors.New("not found")
}
