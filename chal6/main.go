package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	input_reader := bufio.NewReader(os.Stdin)

	input, err := input_reader.ReadString('\n')
	check(err)

	// split by pipe operator
	commands := strings.Split(input, "|")

	var sortOutput bytes.Buffer
	for i, cmd := range commands {
		cmd = strings.TrimSpace(cmd)

		if strings.HasPrefix(cmd, "sort") {
			handleSort(cmd, &sortOutput)
		} else {
			// Handle other commands using os/exec
			handleOtherCommand(cmd, &sortOutput)
		}

		// If there are more commands, prepare the output for piping
		if i < len(commands)-1 {
			sortOutput = *runPipe(&sortOutput)
		}
	}

	// Print the final output
	fmt.Println(sortOutput.String())
}

func handleSort(cmd string, output *bytes.Buffer) {
	fmt.Println("Handle sort called")
	args := strings.Fields(cmd)
	var words []string
	uniq_command := false
	// meaning no -u command (Step 2)
	if len(args) > 2 {
		if args[1] != "-u" {
			fmt.Println("Command should only -u argument, current used:", args[1])
			os.Exit(1)
		}
		is_text_file := isTextFile(args[2])
		if !is_text_file {
			fmt.Println("COmmand should include text file name 2 args, filename:", args[2])
			os.Exit(1)
		}
		uniq_command = true
	} else {
		filename := args[1]
		is_text_file := isTextFile(filename)
		if !is_text_file {
			fmt.Println("COmmand should include text file name 1 args, filename:", filename)
			os.Exit(1)
		}
		file := readFile(filename)

		// new scanner because we need to read lines
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			words = append(words, word)
		}

		// using sort utility
		SelectionSort(words)

		if uniq_command {
			unqiueWords(words)
		}

		for _, word := range words {
			fmt.Fprintln(output, word)
		}
	}
	fmt.Println("Handle sort done")
	err := os.WriteFile("output.txt", output.Bytes(), 0644)
	check(err)

}
