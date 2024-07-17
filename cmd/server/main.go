package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/danielmourad/gochat/models"
)

var (
	clients = make(map[net.Conn]bool)
)

func sendMessageHistory(conn net.Conn) error {
	messages, err := ReadMessages()
	if err != nil {
		return err
	}

	var buffer = make([]byte, 1024)

	for _, msg := range *messages {
		data, err := json.Marshal(&msg)
		if err != nil {
			return err
		}

		copy(buffer, data)
		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Printf("error: sending message: %s\n", err)
		}

		buffer = make([]byte, 1024)
	}

	return nil
}

func handleConnection(conn net.Conn) {
	clients[conn] = true

	defer delete(clients, conn)
	defer conn.Close()

	err := sendMessageHistory(conn)
	if err != nil {
		fmt.Printf("error: sending message history: %s\n", err)
		return
	}

	var buffer = make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("error: receiving message: %s\n", err)
			return
		}
		readData := buffer[:n]

		var msg models.Message
		err = json.Unmarshal(readData, &msg)
		if err != nil {
			fmt.Printf("error: decoding message: %v\n", err)
			continue
		}

		err = WriteMessage(msg)
		if err != nil {
			fmt.Printf("error: saving message: %v\n", err)
		}

		msg.Print()

		for client := range clients {
			if client != conn {
				_, err := client.Write(readData)
				if err != nil {
					fmt.Printf("error: sending message: %s\n", err)
				}
			}
		}

		buffer = make([]byte, 1024)
	}
}

func main() {
	addr := "0.0.0.0"
	port := 4700

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Printf("error: listening: %s\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("gochat server is started\nlistening on %s:%d\n", addr, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error: accepting connection: %s\n", err)
			continue
		}
		go handleConnection(conn)
	}
}
