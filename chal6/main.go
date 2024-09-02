package main

import (
	"bufio"
	"fmt"
)

func main() {
	file := readFile("words.txt")
	file_content := bufio.NewReader(file)
	fmt.Print(file_content)
}
