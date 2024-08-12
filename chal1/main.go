package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter commands (type 'esc' to exit):")

	for {
		// Read input from the user
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)

		// Exit the loop if the user types "esc"
		if input == "esc" {
			fmt.Println("Exiting...")
			break
		}

		// Process the input command
		values := strings.Split(input, " ")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}

		isValidInput := valid(values)
		if !isValidInput {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		fileName := values[2]
		commandType := values[1]
		if commandType == "-c" {
			getSizeInBytes(fileName)
			return
		} else if commandType == "-l" {
			countLines(fileName)
		} else if commandType == "-w" {
			countWords(fileName)
			return
		} else {
			countCharacters(fileName)
			return
		}

	}
}

func check(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

func getSizeInBytes(fileName string) {
	fileStat, err := os.Stat(fileName)
	check(err)

	fileSize := fileStat.Size()
	fmt.Printf("%d %s \n", fileSize, fileName)
}

func readAndReturnFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	check(err)
	return file
}

func countCharacters(fileName string) {
	file := readAndReturnFile(fileName)
	defer file.Close()

	reader := bufio.NewReader(file)
	charCount := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			check(err)
		}
		charCount += utf8.RuneCountInString(line)
	}

	fmt.Printf("%d %s \n", charCount, fileName)

}

func countWords(fileName string) {
	file := readAndReturnFile(fileName)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	// Count the lines.
	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	fmt.Printf("%d\n", count)
}

func countLines(fileName string) {
	file := readAndReturnFile(fileName)
	defer file.Close()

	// bufFileReader := bufio.NewReader(file)
	// lineCount := 0
	// for {
	// 	_, err := bufFileReader.ReadString('\n')
	// 	if err != nil {
	// 		if err.Error() == "EOF" {
	// 			break // End of file is expected; exit the loop
	// 		}
	// 		check(err) // Handle other potential errors
	// 	}
	// 	lineCount++
	// }

	// fmt.Printf("%d %s\n", lineCount, fileName)

	// Alternatively, using bufio.Scanner
	bufFileScanner := bufio.NewScanner(file)
	bufFileScanner.Split(bufio.ScanLines)
	lineCount := 0
	for bufFileScanner.Scan() {
		lineCount++
	}
	if err := bufFileScanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
	fmt.Printf("%d %s\n", lineCount, fileName)
}

func valid(values []string) bool {
	if len(values) > 3 || len(values) < 3 {
		fmt.Println("Wrong input: expected 3 values")
		return false
	}
	if values[0] != "ccwc" {
		fmt.Println("Command not found. Do you mean ccwc?")
		return false
	}
	if values[1] != "-c" && values[1] != "-l" && values[1] != "-w" && values[1] != "-m" {
		fmt.Println("Command not found. Do you mean ccwc -c or -l?")
		return false
	}
	if !(strings.Contains(values[2], ".txt")) {
		fmt.Println("File not found. Do you mean ccwc -c test.txt?")
		return false
	}
	return true
}
