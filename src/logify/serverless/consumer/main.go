package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LogEntry struct {
	Log       string `json:"log"`
	ProjectID string `json:"projectId"`
	Timestamp int64  `json:"timestamp"`
}

var dbPool *pgxpool.Pool

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	var err error
	dbPool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to create database pool: %v\n", err)
	}
}

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	var logs []LogEntry

	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data

		var logEntry LogEntry
		if err := json.Unmarshal(dataBytes, &logEntry); err != nil {
			log.Printf("Error unmarshaling data: %v", err)
			continue
		}

		logs = append(logs, logEntry)
	}

	if len(logs) > 0 {
		if err := bulkInsertLogs(ctx, logs); err != nil {
			return fmt.Errorf("failed to insert logs: %w", err)
		}
	}

	return nil
}

func bulkInsertLogs(ctx context.Context, logs []LogEntry) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, logEntry := range logs {
		batch.Queue("INSERT INTO logs (projectId, log, timestamp) VALUES ($1, $2, $3)",
			logEntry.ProjectID, logEntry.Log, logEntry.Timestamp)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to execute batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
