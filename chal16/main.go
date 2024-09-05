package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

	// Create a channel to handle user input and server responses concurrently
	inputChan := make(chan string)

	// Goroutine to handle server messages
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from server:", err)
				break
			}
			message = strings.TrimSpace(message)
			fmt.Println("Server:", message)

			// Respond to PING messages
			if strings.HasPrefix(message, "PING") {
				response := strings.Replace(message, "PING", "PONG", 1)
				fmt.Fprintf(conn, "%s\r\n", response)
			}

			// Handle JOIN and PART responses (just logging for now)
			if strings.Contains(message, "JOIN") {
				fmt.Println("Joined channel:", extractChannel(message))
			} else if strings.Contains(message, "PART") {
				fmt.Println("Left channel:", extractChannel(message))
			}
		}
	}()

	// Goroutine to handle user input
	go func() {
		userReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			text, _ := userReader.ReadString('\n')
			text = strings.TrimSpace(text)
			inputChan <- text
		}
	}()

	// Use a for-range loop to process user input
	for userInput := range inputChan {
		if strings.HasPrefix(userInput, "/join ") {
			channel := strings.TrimSpace(strings.Split(userInput, " ")[1])
			fmt.Fprintf(conn, "JOIN %s\r\n", channel)
		} else if strings.HasPrefix(userInput, "/part ") {
			channel := strings.TrimSpace(strings.Split(userInput, " ")[1])
			fmt.Fprintf(conn, "PART %s\r\n", channel)
		} else {
			// Other commands can be handled here
			fmt.Fprintf(conn, "%s\r\n", userInput)
		}
	}
}

// extractChannel extracts the channel name from a JOIN or PART message
func extractChannel(message string) string {
	parts := strings.Split(message, " ")
	if len(parts) > 2 {
		return parts[len(parts)-1] // The channel is typically the last part of the message
	}
	return ""
}
