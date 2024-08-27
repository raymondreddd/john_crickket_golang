package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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

	var res []string
	for index, command := range commands {
		snippets := strings.Split(command, " ")

		if index == 0 {
			// check if 1st word is 'cut'
			if snippets[0] != "cut" {
				fmt.Print("Unknown command, use cut")
				os.Exit(1)
			}

			fmt.Print(strings.HasPrefix(snippets[1], "-f"), strings.HasSuffix(snippets[2], ".tsv"), strings.HasSuffix(snippets[2], ".csv"))
			// check for -f
			if len(snippets) < 3 || !strings.HasPrefix(snippets[1], "-f") || !strings.HasSuffix(snippets[2], ".tsv") || !strings.HasSuffix(snippets[2], ".csv") {
				fmt.Print("\n Unknown command, use cut -f1 sample.tsv\n")
				os.Exit(1)
			}

			field := snippets[1][2:]
			fmt.Print("Field number:", field)

			// convert string field to int
			num_field, err := strconv.Atoi(field)
			check(err)

			filename := snippets[2]

			res = extractFields(filename, num_field)
		}

	}

	for _, line := range res {
		fmt.Println(line)
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
		if len(record) < field {
			fmt.Print("This field or column does not exist")
			os.Exit(1)
		}
		res = append(res, record[field])
	}

	return res
}
