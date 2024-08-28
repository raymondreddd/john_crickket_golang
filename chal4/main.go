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
	line_limit := -1
	for index, command := range commands {
		snippets := strings.Split(command, " ")

		if index == 0 {
			// check if 1st word is 'cut'
			if snippets[0] != "cut" {
				fmt.Print("Unknown command, use cut")
				os.Exit(1)
			}

			third_argument := strings.TrimSpace(snippets[2]) // Trim whitespace (including newlines) from the 3rd argument
			filename := ""
			var delimiter string = "" // Default to empty

			if strings.HasPrefix(third_argument, "-d") {
				// validate that file name is present
				if len(snippets) < 4 {
					fmt.Println("\n Icomplete command, use cut -f1 -d, sample.tsv or cut -f1 -d, sample.csv")
					os.Exit(1)
				}

				filename = strings.TrimSpace(snippets[3]) // file name is the 3rd arg
				delimiter = third_argument[2:]            // extract only the dilimeter value
			} else {
				filename = strings.TrimSpace(snippets[2])
			}

			// Check for -f and valid suffixes
			if len(snippets) < 3 || !strings.HasPrefix(snippets[1], "-f") || !(strings.HasSuffix(filename, ".tsv") || strings.HasSuffix(filename, ".csv")) {
				fmt.Println("\nUnknown command, use cut -f1 sample.tsv or cut -f1 sample.csv")
				os.Exit(1)
			}

			field := snippets[1][2:]

			// convert string field to int
			num_field, err := strconv.Atoi(field)
			check(err)

			res = extractFields(filename, num_field, delimiter)
		} else if index == 1 {
			// example: head -n5

			// we dont use [0] because it will contain ""
			first_argument := strings.TrimSpace(snippets[1])

			if first_argument != "head" || len(snippets) < 2 {
				fmt.Println("2nd command should be like head")
				os.Exit(1)
			}

			second_argument := strings.TrimSpace(snippets[2])
			if !strings.HasPrefix(second_argument, "-n") {
				fmt.Println("2nd argument should be like -n4")
				os.Exit(1)
			}

			limit := second_argument[2:]
			// convert string field to int
			num_limit, err := strconv.Atoi(limit)
			check(err)
			line_limit = num_limit
		}

	}

	// cut -f1 -d, fourchords.csv | head -n5
	// cut -f1 four.csv | head -n5
	fmt.Println()
	for line_count, line := range res {
		fmt.Println(line)
		if line_limit != -1 && line_count == line_limit-1 {
			break
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

func extractFields(filename string, field int, delimiter string) []string {
	var res []string

	// Read csv file
	file := readFile(filename)
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Set the delimiter based on the provided value
	if delimiter == "" {
		// Default to tab delimiter for TSV files
		if strings.HasSuffix(filename, ".tsv") {
			reader.Comma = '\t'
		} else {
			reader.Comma = ',' // Default to comma for CSV or other files
		}
	} else {
		reader.Comma = rune(delimiter[0])
	}

	// Enable lazy quotes to handle improperly formatted quoted fields
	reader.LazyQuotes = true

	// Read all records
	records, err := reader.ReadAll()
	check(err)

	fmt.Println("delimei:", delimiter)
	// Extract the specified field from each record
	for _, record := range records {
		if len(record) < field {
			fmt.Print("This field or column does not exist")
			os.Exit(1)
		}
		// cut -f1 -d, four.csv | head -n5
		// cut -f1 four.csv | head -n5
		// Append the field's value directly to the result
		res = append(res, record[field-1])
	}

	return res
}
