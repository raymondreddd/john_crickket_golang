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
