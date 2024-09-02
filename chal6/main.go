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
	args := strings.Fields(cmd)

	// meaning no -u command (Step 2)
	if len(args) < 2 {
		if args[0] != "-u" {
			fmt.Print("Command should only -u argument")
			os.Exit(1)
		}
		is_text_file := isTextFile(args[1])
		if !is_text_file {
			fmt.Print("COmmand should include text file name")
			os.Exit(1)
		}
	} else {
		filename := args[0]
		is_text_file := isTextFile(filename)
		if !is_text_file {
			fmt.Print("COmmand should include text file name")
			os.Exit(1)
		}
		file := readFile(filename)

		// new scanner because we need to read lines
		file_content := bufio.NewScanner(file)
		fmt.Print(file_content)

	}
	// Here you implement the sort logic
	// For simplicity, we can assume it reads from stdin and writes to the output buffer
	fmt.Fprintln(output, "Sorted data would go here.")
}
