package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func readFile(filename string) *os.File {
	file, err := os.Open(filename)
	check(err)
	return file
}

func check(err error) {
	if err != nil {
		fmt.Println("Error :", err)
		os.Exit(1)
	}
}

// dont pass pointer
func SelectionSort(arr []string) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		min_index := i
		// unsorted part of the array
		for j := i + 1; j < n; j++ {
			// update min_index
			if arr[j] < arr[min_index] {
				min_index = j
			}
		}
		// Swap the minimum element with the curr elemnt
		arr[i], arr[min_index] = arr[min_index], arr[i]
	}
}

func unqiueWords(words []string) []string {
	if len(words) == 0 {
		return words
	}

	uniqueWords := []string{words[0]}
	for i := 1; i < len(words); i++ {
		if words[i] != words[i-1] {
			uniqueWords = append(uniqueWords, words[i])
		}
	}
	return uniqueWords
}

func runPipe(input *bytes.Buffer) *bytes.Buffer {
	output := new(bytes.Buffer)
	output.Write(input.Bytes())
	return output
}

func isTextFile(input string) bool {
	return strings.HasSuffix(input, ".txt")
}

func handleOtherCommand(cmd string, input *bytes.Buffer) {
	// for example `head -n 5` to []string{"head", "-n", "5"}
	args := strings.Fields(cmd)
	name := args[0]
	cmdArgs := args[1:]

	command := exec.Command(name, cmdArgs...)

	// redirect the input for command to come from `input` buffer or the sorted worsd
	// which contains the output of the above command
	command.Stdin = input

	// to store the output
	var out bytes.Buffer

	// redirect the command output to out buffer
	command.Stdout = &out

	err := command.Run()
	check(err)

	// Clears the original input buffer.
	input.Reset()

	// this is chaining, output of (out) into `input` buffer
	// e.g. next input can be uniq, so it will place only uniq words,etc.
	input.Write(out.Bytes())
}
