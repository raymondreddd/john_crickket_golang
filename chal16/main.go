package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// server conatins host and port, meaning address we connect to
	server := "irc.freenode.net:6667"
	nick := "Muku"
	user := "guest 0 * :Challenge 16 go Client"

	// Connecting to irc server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send NICK and USER commands
	fmt.Fprintf(conn, "NICK %s\r\n", nick)
	fmt.Fprintf(conn, "USER %s\r\n", user)

	// Read and print server responses
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			break
		}
		fmt.Print(message)
	}
}
