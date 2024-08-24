package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
COmperession tool

1. Create hashmap of frequencies of each char/rune
*/

func main() {
	file := readFile("test.txt")

	genFrequencyMap(file)
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func readFile(file_name string) *os.File {
	file, err := os.Open(file_name)
	check(err)
	return file
}

func genFrequencyMap(file *os.File) {
	reader := bufio.NewReader(file)
	defer file.Close()

	freq_map := make(map[rune]int)

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading file:", err)
			return
		}
		freq_map[char]++
	}

	// ch is character and frq is frequency
	for ch, frq := range freq_map {
		fmt.Printf("Character: %c %d \n", ch, frq)
	}

}
