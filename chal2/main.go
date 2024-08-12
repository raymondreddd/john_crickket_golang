package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ParanStack struct {
	elements []rune
}

// Push adds an element to the ParanStack.
func (s *ParanStack) Push(element rune) {
	s.elements = append(s.elements, element)
}

// Pop removes and returns the top element of the ParanStack. Returns an error if the ParanStack is empty.
func (s *ParanStack) Pop() (rune, error) {
	if len(s.elements) == 0 {
		return 0, errors.New("ParanStack is empty")
	}

	// Get the last element
	topElement := s.elements[len(s.elements)-1]

	// Remove the last element
	s.elements = s.elements[:len(s.elements)-1]

	return topElement, nil
}

// To check general Error
func checkError(e error) {
	if e != nil {
		fmt.Println("error", e)
		os.Exit(1)
	}
}

// fucntion to open file
func readFile(file_name string) *os.File {
	file, err := os.Open(file_name)
	checkError(err)
	return file
}
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		checkError(err)

		values := strings.Split(input, " ")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
			if values[i] == "esc" {
				os.Exit(1)
			}
		}
		res1 := step1()
		if !res1 {
			os.Exit(0)
		}
	}
}

// Step 1 divide
func step1() bool {
	file := readFile("./tests/step1/valid.json")
	defer file.Close()
	reader := bufio.NewReader(file)
	valid_braces := ParanStack{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			fmt.Println("Error reading file:")
			checkError(err)
		}
		for _, ch := range line {
			if ch == '{' {
				valid_braces.Push(ch)
			} else if ch == '}' {
				valid_braces.Pop()
			}
		}
	}
	if (len(valid_braces.elements)) == 0 {
		fmt.Println("Valid JSON")
		return true
	} else {
		fmt.Println("Invalid JSON")
		return false
	}
}
