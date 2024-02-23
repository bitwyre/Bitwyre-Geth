package utils

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// TransactionLog represents the structure of the log to be saved
type TransactionLog struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

func RecordTransaction(tx *types.Transaction, txType string, filename string) error {
	// Prepare the data to be logged
	data := TransactionLog{
		Hash:      tx.Hash().Hex(),
		Timestamp: time.Now(),
		Type:      txType,
	}

	// Marshal data to JSON format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Error("Failed to open file")
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Error("Failed to write to file")
		return err
	}

	// Write a newline for readability if logging multiple entries
	_, err = file.WriteString("\n")
	return err
}
