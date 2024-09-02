package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	file := readFile("words.txt")
	file_content := bufio.NewReader(file)
	fmt.Print(file_content)

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
	// Here you implement the sort logic
	// For simplicity, we can assume it reads from stdin and writes to the output buffer
	fmt.Fprintln(output, "Sorted data would go here.")
}
