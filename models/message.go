package models

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Text   Text   `json:"text"`
}

func NewMessage(author string, text string) (*Message, error) {
	if len(text) > 256 {
		return nil, fmt.Errorf("text can not be more than 256 characters")
	}

	return &Message{
		Author: author,
		Text:   Text(text),
	}, nil
}

func (msg *Message) Print() {
	fmt.Println(msg.Author, "@", msg.Date, "->", msg.Text.String())
}

func (msg *Message) Send(conn net.Conn) error {
	msg.Date = time.Now().Format("2006-01-02 15:04:05")

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
