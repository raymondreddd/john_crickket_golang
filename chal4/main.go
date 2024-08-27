package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {

	filename := "sample.tsv"
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
