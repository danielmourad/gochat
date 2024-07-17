package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/danielmourad/gochat/models"
)

const dbFileName string = "database.json"

func init() {
	err := initializeDatabase()
	if err != nil {
		fmt.Printf("error: unable to initialize database: %s\n", err)
		os.Exit(1)
	}
}

type Database struct {
	Messages []models.Message `json:"messages"`
}

func initializeDatabase() error {
	_, err := os.Stat(dbFileName)
	if os.IsNotExist(err) {
		dbFile, err := os.Create(dbFileName)
		if err != nil {
			return err
		}
		defer dbFile.Close()

		var db Database
		jsonData, err := json.MarshalIndent(db, "", "  ")
		if err != nil {
			return err
		}

		_, err = dbFile.Write(jsonData)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func ReadMessages() (*[]models.Message, error) {
	data, err := os.ReadFile(dbFileName)
	if err != nil {
		return nil, err
	}

	var db Database
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	return &db.Messages, nil
}

func WriteMessage(msg models.Message) error {
	dbFile, err := os.OpenFile(dbFileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer dbFile.Close()

	data, err := io.ReadAll(dbFile)
	if err != nil {
		return err
	}

	var db Database
	err = json.Unmarshal(data, &db)
	if err != nil {
		return err
	}

	db.Messages = append(db.Messages, msg)

	jsonData, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	_, err = dbFile.WriteAt(jsonData, 0)
	if err != nil {
		return err
	}

	return nil
}
