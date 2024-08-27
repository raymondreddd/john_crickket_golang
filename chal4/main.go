package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {

	// Taking user input
	input_reader := bufio.NewReader(os.Stdin)

	input, err := input_reader.ReadString('\n')
	check(err)

	// split into multiple user commands by pipe operator
	commands := strings.Split(input, "|")

	if len(commands) == 0 {
		return
	}

	for index, command := range commands {
		snippets := strings.Split(command, " ")

		if index == 0 {
			// check if 1st word is 'cut'
			if snippets[0] != "cut" {
				fmt.Print("Unknown command, use cut")
				os.Exit(1)
			}

			// check for -f
			if len(snippets) < 3 || !strings.HasPrefix(snippets[1], "-f") {
				fmt.Print("Unknown command, use cut -f1 sample.tsv")
				os.Exit(1)
			}

			field := snippets[1][2:]
			fmt.Print("Field number:", field)

		}

	}

}

func readFile(filename string) *os.File {
	file, err := os.Open(filename)
	check(err)

	return file
}

func check(err error) {
	if err != nil {
		fmt.Print("error:", err)
		os.Exit(1)
	}
}

func extractFields(filename string, field int) []string {
	var res []string

	// Read csv file
	file := readFile(filename)
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Check the file extension for tsv, if yes then set the delimiter to a tab for TSV files
	if strings.HasSuffix(filename, ".tsv") {
		reader.Comma = '\t'
	}

	records, err := reader.ReadAll()
	check(err)

	for _, record := range records {
		fmt.Println(record)
	}

	return res
}
