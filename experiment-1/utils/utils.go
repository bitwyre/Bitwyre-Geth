package utils

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionLog represents the structure of the log to be saved
type TransactionLog struct {
	Hash      string    `json:"hash"`
	NodeId    string    `json:"node_id"`
	Timestamp time.Time `json:"timestamp"`
}

func RecordTransaction(tx *types.Transaction, nodeId string, filename string) error {
	// Prepare the data to be logged
	data := TransactionLog{
		Hash:      tx.Hash().Hex(),
		NodeId:    nodeId,
		Timestamp: time.Now(),
	}

	// Marshal data to JSON format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	// Write a newline for readability if logging multiple entries
	_, err = file.WriteString("\n")
	return err
}
