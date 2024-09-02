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

func handleOtherCommand(cmd string, input *bytes.Buffer) {
	args := strings.Fields(cmd)
	name := args[0]
	cmdArgs := args[1:]

	command := exec.Command(name, cmdArgs...)
	command.Stdin = input
	var out bytes.Buffer
	command.Stdout = &out

	err := command.Run()
	check(err)

	// Copy output to the main buffer
	input.Reset()
	input.Write(out.Bytes())
}
