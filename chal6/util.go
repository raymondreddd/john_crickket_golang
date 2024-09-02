package main

import (
	"fmt"
	"os"
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
