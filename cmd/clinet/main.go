package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/danielmourad/gochat/models"
)

func receiveMessage(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("error: receiving message: %s\n", err)
			return
		}
		readData := buffer[:n]
		trimmed := bytes.TrimRight(readData, "\x00")

		var received models.Message
		err = json.Unmarshal(trimmed, &received)
		if err != nil {
			fmt.Printf("error: decoding message: %v\n", err)
			continue
		}

		fmt.Print("\r" + " " + "\r")
		received.Print()
		fmt.Print("-> ")

		buffer = make([]byte, 1024)
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("error: you need to provide name when running the client\n")
		os.Exit(1)
	}

	author := os.Args[1]
	serverAddr := "127.0.0.1"
	serverPort := 4700

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverAddr, serverPort))
	if err != nil {
		fmt.Printf("error: connecting to server: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("connection with %s:%d gochat server established\n", serverAddr, serverPort)
	fmt.Printf("your display name is: %s\n", author)

	go receiveMessage(conn)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("-> ")
		scanner.Scan()

		text := scanner.Text()

		msg, err := models.NewMessage(author, text)
		if err != nil {
			fmt.Printf("error: creating message: %s\n", err)
			continue
		}
		if msg.Text.String() == "exit" {
			fmt.Println("bye!")
			os.Exit(0)
		}

		err = msg.Send(conn)
		if err != nil {
			fmt.Printf("error: sending message: %s\n", err)
		}

		fmt.Print("\033[1A\033[K")
		msg.Print()
	}
}
