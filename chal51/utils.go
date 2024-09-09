package main

import (
	"bufio"
	"os"
	"strings"
)

// Load forbidden hosts from file into a set (map with empty struct values)
func loadBannedHosts(filePath string) error {
	BannedHosts = make(map[string]struct{})

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			BannedHosts[line] = struct{}{}
		}
	}
	return scanner.Err()
}
