package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"encoding/csv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	client "github.com/influxdata/influxdb1-client/v2"
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

func SaveToInflux(tx *types.Transaction) {

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("GETH_METRICS_INFLUXDB_ENDPOINT"),
		Username: os.Getenv("GETH_METRICS_INFLUXDB_USERNAME"),
		Password: os.Getenv("GETH_METRICS_INFLUXDB_PASSWORD"),
		Timeout:  10 * time.Second,
	})
	if err != nil {
		log.Info("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	bp, err2 := client.NewBatchPoints(client.BatchPointsConfig{
		Database: os.Getenv("GETH_METRICS_INFLUXDB_DATABASE"),
		Precision: "s",
	})

	if err2 != nil {
		log.Error("error creating batch points: %w", err2)
		return
	}

	txHash := tx.Hash().Hex()
	currentTime := time.Now()

	log.Info("Attempt to write transaction with hash: ", txHash, nil)

	tags := map[string]string{"transaction_hash": txHash}
	fields := map[string]interface{}{
		"recorded_at": currentTime.Unix(),
	}

	pt, err := client.NewPoint("transaction_metrics", tags, fields, currentTime)
	if err != nil {
		log.Error("error creating new point for InfluxDB", "err", err)
	}
	log.Info("The data is send")
	bp.AddPoint(pt)

	if err := c.Write(bp); err != nil {
		log.Error("Failed to write batch points to InfluxDB", "err", err)
	}

	q := client.NewQuery("SELECT * FROM transaction_metrics", os.Getenv("GETH_METRICS_INFLUXDB_DATABASE"), "")
	
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		file, err := os.Create("output.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"timestamp", "recorded_at", "tx_id"})

		for _, res := range response.Results {
			for _, series := range res.Series {
				for _, row := range series.Values {
					if len(row) < 2 {
						continue
					}
					timeVal, ok := row[0].(string)
					if !ok {
						fmt.Println("Error: Unable to convert _time to string")
						continue
					}
					value, ok := row[1].(json.Number)
					if !ok {
						fmt.Println("Error: Unable to convert _value to string")
						continue
					}

					value2, ok := row[2].(string)
					if !ok {
						fmt.Println("Error: Unable to convert _value to string")
						continue
					}
					writer.Write([]string{timeVal, string(value), string(value2)}) // Modify the row based on your data schema
				}
			}
		}
	}

}
