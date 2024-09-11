package main

import (
	"bufio"
	"fmt"
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

func loadBannedWords(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error loading banned words: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if word != "" {
			BannedWords[word] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading banned words: %v\n", err)
	}
}
